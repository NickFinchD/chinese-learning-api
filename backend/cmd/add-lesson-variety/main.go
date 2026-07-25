// Command add-lesson-variety enriches every existing lesson (in any course)
// with two things, without removing or reordering any of its existing
// steps:
//
//   - A "переводится → слово" (translation_to_word) quiz right after every
//     existing "Как переводится X?" (word_to_translation) quiz, testing the
//     same word in reverse — mirroring how each word already gets a 'word'
//     step plus a forward quiz.
//   - One sentence_builder step, for lessons that don't already have one,
//     inserted right after the lesson's vocabulary/grammar block (or at the
//     start, for review-only lessons with no 'word'/'grammar' steps).
//     Reused from the shared sentence_exercises pool for the course's HSK
//     level — the pool is far smaller than the lesson count, so exercises
//     necessarily repeat across lessons.
//
//	go run ./cmd/add-lesson-variety
//
// Idempotent: reruns skip words that already have a reverse quiz and
// lessons that already have a sentence_builder step.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	directionWordToTranslation = "word_to_translation"
	directionTranslationToWord = "translation_to_word"
)

type wordInfo struct {
	ID          int64
	Hanzi       string
	Pinyin      string
	Translation string
	HSKLevel    int16
}

type quizInfo struct {
	ID        int64
	Direction string
	Hanzi     *string
}

type stepRow struct {
	ID        int64
	LessonID  int64
	StepType  string
	EntityID  int64
	SortOrder int
}

type newStep struct {
	StepType string
	EntityID int64
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	words, err := loadWords(ctx, db)
	if err != nil {
		log.Fatalf("failed to load words: %v", err)
	}

	wordByHanzi := make(map[string]wordInfo, len(words))
	wordIDsByLevel := map[int16][]int64{}
	wordByID := make(map[int64]wordInfo, len(words))

	for _, w := range words {
		wordByHanzi[w.Hanzi] = w
		wordByID[w.ID] = w
		wordIDsByLevel[w.HSKLevel] = append(wordIDsByLevel[w.HSKLevel], w.ID)
	}

	quizzes, err := loadQuizzes(ctx, db)
	if err != nil {
		log.Fatalf("failed to load quizzes: %v", err)
	}

	quizByID := make(map[int64]quizInfo, len(quizzes))
	reverseQuizByHanzi := make(map[string]int64)

	for _, q := range quizzes {
		quizByID[q.ID] = q
		if q.Direction == directionTranslationToWord && q.Hanzi != nil {
			reverseQuizByHanzi[*q.Hanzi] = q.ID
		}
	}

	lessonHSK, err := loadLessonHSKLevels(ctx, db)
	if err != nil {
		log.Fatalf("failed to load lesson HSK levels: %v", err)
	}

	sentencePool, err := loadSentencePoolByHSK(ctx, db)
	if err != nil {
		log.Fatalf("failed to load sentence exercises: %v", err)
	}

	rng := rand.New(rand.NewSource(7))

	for level, pool := range sentencePool {
		rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		sentencePool[level] = pool
	}

	sentenceCursor := map[int16]int{}

	steps, err := loadSteps(ctx, db)
	if err != nil {
		log.Fatalf("failed to load lesson steps: %v", err)
	}

	byLesson := map[int64][]stepRow{}
	for _, s := range steps {
		byLesson[s.LessonID] = append(byLesson[s.LessonID], s)
	}

	lessonIDs := make([]int64, 0, len(byLesson))
	for id := range byLesson {
		lessonIDs = append(lessonIDs, id)
	}
	sort.Slice(lessonIDs, func(i, j int) bool { return lessonIDs[i] < lessonIDs[j] })

	lessonsChanged := 0
	reverseQuizzesAdded := 0
	sentenceStepsAdded := 0

	for _, lessonID := range lessonIDs {

		original := byLesson[lessonID]
		sort.Slice(original, func(i, j int) bool { return original[i].SortOrder < original[j].SortOrder })

		hasSentenceBuilder := false
		lastVocabIdx := -1

		for i, s := range original {
			if s.StepType == "sentence_builder" {
				hasSentenceBuilder = true
			}
			if s.StepType == "word" || s.StepType == "grammar" {
				lastVocabIdx = i
			}
		}

		sentenceInsertAt := lastVocabIdx + 1 // 0 if no word/grammar steps

		hskLevel := lessonHSK[lessonID]
		pool := sentencePool[hskLevel]

		planned := make([]newStep, 0, len(original)+len(original)/2+1)
		changed := false

		for i, s := range original {

			if !hasSentenceBuilder && i == sentenceInsertAt && len(pool) > 0 {
				cursor := sentenceCursor[hskLevel]
				exerciseID := pool[cursor%len(pool)]
				sentenceCursor[hskLevel] = cursor + 1

				planned = append(planned, newStep{StepType: "sentence_builder", EntityID: exerciseID})
				sentenceStepsAdded++
				changed = true
			}

			planned = append(planned, newStep{StepType: s.StepType, EntityID: s.EntityID})

			if s.StepType != "quiz" {
				continue
			}

			quiz, ok := quizByID[s.EntityID]
			if !ok || quiz.Direction != directionWordToTranslation || quiz.Hanzi == nil {
				continue
			}

			w, ok := wordByHanzi[*quiz.Hanzi]
			if !ok {
				continue
			}

			reverseID, ok := reverseQuizByHanzi[w.Hanzi]
			if !ok {

				distractors := pickDistractorWords(rng, wordIDsByLevel[w.HSKLevel], wordByID, w.ID, 3)
				if len(distractors) < 3 {
					continue
				}

				var err error
				reverseID, err = createReverseQuiz(ctx, db, w, distractors, rng)
				if err != nil {
					log.Fatalf("failed to create reverse quiz for %q: %v", w.Hanzi, err)
				}

				reverseQuizByHanzi[w.Hanzi] = reverseID
			}

			planned = append(planned, newStep{StepType: "quiz", EntityID: reverseID})
			reverseQuizzesAdded++
			changed = true
		}

		if !hasSentenceBuilder && sentenceInsertAt == len(original) && len(pool) > 0 {
			cursor := sentenceCursor[hskLevel]
			exerciseID := pool[cursor%len(pool)]
			sentenceCursor[hskLevel] = cursor + 1

			planned = append(planned, newStep{StepType: "sentence_builder", EntityID: exerciseID})
			sentenceStepsAdded++
			changed = true
		}

		if !changed {
			continue
		}

		if err := replaceLessonSteps(ctx, db, lessonID, planned); err != nil {
			log.Fatalf("failed to update lesson %d: %v", lessonID, err)
		}

		lessonsChanged++
	}

	fmt.Printf(
		"Done. %d lesson(s) updated, %d reverse quiz step(s) added, %d sentence-builder step(s) added.\n",
		lessonsChanged, reverseQuizzesAdded, sentenceStepsAdded,
	)
}

func loadWords(ctx context.Context, db *pgxpool.Pool) ([]wordInfo, error) {
	rows, err := db.Query(ctx, `SELECT id, hanzi, pinyin, translation, hsk_level FROM words ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]wordInfo, 0)
	for rows.Next() {
		var w wordInfo
		if err := rows.Scan(&w.ID, &w.Hanzi, &w.Pinyin, &w.Translation, &w.HSKLevel); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

func loadQuizzes(ctx context.Context, db *pgxpool.Pool) ([]quizInfo, error) {
	rows, err := db.Query(ctx, `SELECT id, direction, hanzi FROM quizzes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quizzes := make([]quizInfo, 0)
	for rows.Next() {
		var q quizInfo
		if err := rows.Scan(&q.ID, &q.Direction, &q.Hanzi); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, rows.Err()
}

func loadLessonHSKLevels(ctx context.Context, db *pgxpool.Pool) (map[int64]int16, error) {
	rows, err := db.Query(ctx, `
		SELECT l.id, c.hsk_level
		FROM lessons l
		JOIN courses c ON c.id = l.course_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]int16)
	for rows.Next() {
		var id int64
		var level int16
		if err := rows.Scan(&id, &level); err != nil {
			return nil, err
		}
		result[id] = level
	}
	return result, rows.Err()
}

func loadSentencePoolByHSK(ctx context.Context, db *pgxpool.Pool) (map[int16][]int64, error) {
	rows, err := db.Query(ctx, `SELECT id, hsk_level FROM sentence_exercises ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int16][]int64)
	for rows.Next() {
		var id int64
		var level int16
		if err := rows.Scan(&id, &level); err != nil {
			return nil, err
		}
		result[level] = append(result[level], id)
	}
	return result, rows.Err()
}

func loadSteps(ctx context.Context, db *pgxpool.Pool) ([]stepRow, error) {
	rows, err := db.Query(ctx, `
		SELECT id, lesson_id, step_type, entity_id, sort_order
		FROM lesson_steps
		ORDER BY lesson_id, sort_order
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]stepRow, 0)
	for rows.Next() {
		var s stepRow
		if err := rows.Scan(&s.ID, &s.LessonID, &s.StepType, &s.EntityID, &s.SortOrder); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

// pickDistractorWords returns up to n words distinct by hanzi from the
// given id pool, excluding excludeID.
func pickDistractorWords(rng *rand.Rand, ids []int64, byID map[int64]wordInfo, excludeID int64, n int) []wordInfo {

	shuffled := append([]int64{}, ids...)
	rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	excludeHanzi := byID[excludeID].Hanzi
	seen := map[string]bool{excludeHanzi: true}
	result := make([]wordInfo, 0, n)

	for _, id := range shuffled {
		if id == excludeID {
			continue
		}
		w := byID[id]
		if seen[w.Hanzi] {
			continue
		}
		seen[w.Hanzi] = true
		result = append(result, w)
		if len(result) == n {
			break
		}
	}

	return result
}

func createReverseQuiz(ctx context.Context, db *pgxpool.Pool, w wordInfo, distractors []wordInfo, rng *rand.Rand) (int64, error) {

	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	question := fmt.Sprintf("Как по-китайски «%s»?", w.Translation)

	var quizID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO quizzes (question, hsk_level, direction, hanzi, pinyin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, question, w.HSKLevel, directionTranslationToWord, w.Hanzi, w.Pinyin).Scan(&quizID); err != nil {
		return 0, err
	}

	type opt struct {
		Hanzi     string
		Pinyin    string
		IsCorrect bool
	}

	options := make([]opt, 0, len(distractors)+1)
	options = append(options, opt{Hanzi: w.Hanzi, Pinyin: w.Pinyin, IsCorrect: true})
	for _, d := range distractors {
		options = append(options, opt{Hanzi: d.Hanzi, Pinyin: d.Pinyin, IsCorrect: false})
	}

	rng.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })

	for idx, o := range options {
		if _, err := tx.Exec(ctx, `
			INSERT INTO quiz_options (quiz_id, option_text, pinyin, is_correct, sort_order)
			VALUES ($1, $2, $3, $4, $5)
		`, quizID, o.Hanzi, o.Pinyin, o.IsCorrect, idx+1); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return quizID, nil
}

func replaceLessonSteps(ctx context.Context, db *pgxpool.Pool, lessonID int64, planned []newStep) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM lesson_steps WHERE lesson_id = $1`, lessonID); err != nil {
		return err
	}

	for i, s := range planned {
		if _, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, $2, $3, $4)
		`, lessonID, s.StepType, s.EntityID, i+1); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
