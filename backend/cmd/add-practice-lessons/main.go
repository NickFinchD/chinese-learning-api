// Command add-practice-lessons inserts a dedicated practice lesson right
// after every existing lesson of the "HSK 1" course (except the last one,
// which has nothing left to prepare the learner for), made up entirely of
// quiz steps — both directions where available — for the words that lesson
// just introduced as flashcards.
//
// Existing lessons are renumbered onto odd lesson_number slots (1, 3, 5, ...)
// to make room for the new practice lessons at the even slots between them
// (2, 4, 6, ...), so every "lesson N followed by lesson N+1" relationship the
// frontend relies on (LessonPage's "next lesson" lookup is a strict
// lesson_number+1 match) keeps working with no gaps.
//
// For a lesson that introduced no new words (a pure grammar/review lesson),
// the practice lesson instead just repeats that lesson's own quiz steps
// as-is, so every original lesson still gets a practice lesson after it.
//
// Quizzes are looked up by hanzi + direction and reused when they already
// exist (most words already have both a word_to_translation and a
// translation_to_word quiz somewhere in the DB, from cmd/seed and
// cmd/add-lesson-variety); missing ones are generated with the same
// distractor-picking approach those two commands use.
//
// Run once:
//
//	go run ./cmd/add-practice-lessons
//
// Not idempotent: it aborts immediately if the course's lesson_number
// sequence isn't contiguous (1, 2, 3, ...), which is what it looks like after
// a first successful run, to avoid renumbering or inserting a second time.
//
// Because every lesson's lesson_number shifts, it also clears
// user_lesson_progress and user_course_progress for this course at the end
// (lesson_id references stay valid — nothing else needs to change — but a
// user's in-progress position, expressed as "up to lesson_number X", would
// otherwise silently point at the wrong content).
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
	courseTitle                = "HSK 1"
)

type wordInfo struct {
	ID          int64
	Hanzi       string
	Pinyin      string
	Translation string
	HSKLevel    int16
}

type lessonRow struct {
	ID        int64
	OldNumber int
	Title     string
}

type stepRow struct {
	LessonID  int64
	StepType  string
	EntityID  int64
	SortOrder int
}

type plannedStep struct {
	StepType string
	EntityID int64
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	var courseID int64
	if err := db.QueryRow(ctx, `SELECT id FROM courses WHERE title = $1`, courseTitle).Scan(&courseID); err != nil {
		log.Fatalf("failed to find course %q: %v", courseTitle, err)
	}

	lessons, err := loadLessons(ctx, db, courseID)
	if err != nil {
		log.Fatalf("failed to load lessons: %v", err)
	}

	if len(lessons) < 2 {
		log.Fatalf("course %q has only %d lesson(s), nothing to do", courseTitle, len(lessons))
	}

	for i := 1; i < len(lessons); i++ {
		if lessons[i].OldNumber != lessons[i-1].OldNumber+1 {
			log.Fatalf(
				"lesson_number sequence for %q is not contiguous (%d follows %d) - looks already processed, aborting",
				courseTitle, lessons[i].OldNumber, lessons[i-1].OldNumber,
			)
		}
	}

	words, err := loadWords(ctx, db)
	if err != nil {
		log.Fatalf("failed to load words: %v", err)
	}

	wordByID := make(map[int64]wordInfo, len(words))
	translationByID := make(map[int64]string, len(words))
	wordIDsByLevel := map[int16][]int64{}
	allWordIDs := make([]int64, 0, len(words))

	for _, w := range words {
		wordByID[w.ID] = w
		translationByID[w.ID] = w.Translation
		wordIDsByLevel[w.HSKLevel] = append(wordIDsByLevel[w.HSKLevel], w.ID)
		allWordIDs = append(allWordIDs, w.ID)
	}

	forwardQuizByHanzi, reverseQuizByHanzi, err := loadQuizzesByHanzi(ctx, db)
	if err != nil {
		log.Fatalf("failed to load quizzes: %v", err)
	}

	steps, err := loadSteps(ctx, db, lessons)
	if err != nil {
		log.Fatalf("failed to load lesson steps: %v", err)
	}

	stepsByLesson := map[int64][]stepRow{}
	for _, s := range steps {
		stepsByLesson[s.LessonID] = append(stepsByLesson[s.LessonID], s)
	}
	for id := range stepsByLesson {
		sort.Slice(stepsByLesson[id], func(i, j int) bool {
			return stepsByLesson[id][i].SortOrder < stepsByLesson[id][j].SortOrder
		})
	}

	rng := rand.New(rand.NewSource(20260731))

	// Renumber existing lessons onto odd slots first - old numbers only ever
	// grow when doubled-minus-one, so no two lessons collide mid-update even
	// without a unique constraint to enforce it.
	for _, l := range lessons {
		newNumber := 2*l.OldNumber - 1
		if _, err := db.Exec(ctx, `UPDATE lessons SET lesson_number = $1 WHERE id = $2`, newNumber, l.ID); err != nil {
			log.Fatalf("failed to renumber lesson %d: %v", l.ID, err)
		}
	}

	practiceLessonsAdded := 0
	quizzesCreated := 0

	for i := 0; i < len(lessons)-1; i++ {

		lesson := lessons[i]
		lessonSteps := stepsByLesson[lesson.ID]

		var newWords []wordInfo
		for _, s := range lessonSteps {
			if s.StepType == "word" {
				if w, ok := wordByID[s.EntityID]; ok {
					newWords = append(newWords, w)
				}
			}
		}

		var planned []plannedStep

		if len(newWords) == 0 {
			// Nothing new was taught (grammar/review lesson) - fall back to
			// repeating this lesson's own quiz steps as the practice content.
			for _, s := range lessonSteps {
				if s.StepType == "quiz" {
					planned = append(planned, plannedStep{StepType: "quiz", EntityID: s.EntityID})
				}
			}
		} else {
			for _, w := range newWords {

				forwardID, ok := forwardQuizByHanzi[w.Hanzi]
				if !ok {
					distractors := pickDistractors(rng, allWordIDs, translationByID, w.ID, 3)
					if len(distractors) == 3 {
						id, err := createForwardQuiz(ctx, db, w, distractors, rng)
						if err != nil {
							log.Fatalf("failed to create forward quiz for %q: %v", w.Hanzi, err)
						}
						forwardID = id
						forwardQuizByHanzi[w.Hanzi] = id
						quizzesCreated++
						ok = true
					}
				}
				if ok {
					planned = append(planned, plannedStep{StepType: "quiz", EntityID: forwardID})
				}

				reverseID, ok := reverseQuizByHanzi[w.Hanzi]
				if !ok {
					distractors := pickDistractorWords(rng, wordIDsByLevel[w.HSKLevel], wordByID, w.ID, 3)
					if len(distractors) == 3 {
						id, err := createReverseQuiz(ctx, db, w, distractors, rng)
						if err != nil {
							log.Fatalf("failed to create reverse quiz for %q: %v", w.Hanzi, err)
						}
						reverseID = id
						reverseQuizByHanzi[w.Hanzi] = id
						quizzesCreated++
						ok = true
					}
				}
				if ok {
					planned = append(planned, plannedStep{StepType: "quiz", EntityID: reverseID})
				}
			}

			rng.Shuffle(len(planned), func(a, b int) { planned[a], planned[b] = planned[b], planned[a] })
		}

		if len(planned) == 0 {
			log.Printf("skipping practice lesson after %q: no quiz content available", lesson.Title)
			continue
		}

		practiceNumber := 2 * lesson.OldNumber
		title := fmt.Sprintf("Практика: %s", lesson.Title)
		description := fmt.Sprintf("Закрепите слова из урока «%s» с помощью тестов.", lesson.Title)

		if err := createPracticeLesson(ctx, db, courseID, practiceNumber, title, description, planned); err != nil {
			log.Fatalf("failed to create practice lesson after %q: %v", lesson.Title, err)
		}

		practiceLessonsAdded++
	}

	if _, err := db.Exec(ctx, `
		DELETE FROM user_lesson_progress
		WHERE lesson_id IN (SELECT id FROM lessons WHERE course_id = $1)
	`, courseID); err != nil {
		log.Fatalf("failed to reset user_lesson_progress: %v", err)
	}

	if _, err := db.Exec(ctx, `DELETE FROM user_course_progress WHERE course_id = $1`, courseID); err != nil {
		log.Fatalf("failed to reset user_course_progress: %v", err)
	}

	fmt.Printf(
		"Done. %d practice lesson(s) added, %d new quiz(zes) created, progress reset for course %q.\n",
		practiceLessonsAdded, quizzesCreated, courseTitle,
	)
}

func loadLessons(ctx context.Context, db *pgxpool.Pool, courseID int64) ([]lessonRow, error) {
	rows, err := db.Query(ctx, `
		SELECT id, lesson_number, title
		FROM lessons
		WHERE course_id = $1
		ORDER BY lesson_number
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]lessonRow, 0)
	for rows.Next() {
		var l lessonRow
		if err := rows.Scan(&l.ID, &l.OldNumber, &l.Title); err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, rows.Err()
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

// loadQuizzesByHanzi maps a word's hanzi to an existing quiz id, separately
// per direction, picking the lowest quiz id when more than one matches so
// results are deterministic across reruns.
func loadQuizzesByHanzi(ctx context.Context, db *pgxpool.Pool) (map[string]int64, map[string]int64, error) {
	rows, err := db.Query(ctx, `
		SELECT id, direction, hanzi
		FROM quizzes
		WHERE hanzi IS NOT NULL
		ORDER BY id
	`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	forward := map[string]int64{}
	reverse := map[string]int64{}

	for rows.Next() {
		var id int64
		var direction string
		var hanzi string
		if err := rows.Scan(&id, &direction, &hanzi); err != nil {
			return nil, nil, err
		}

		switch direction {
		case directionWordToTranslation:
			if _, ok := forward[hanzi]; !ok {
				forward[hanzi] = id
			}
		case directionTranslationToWord:
			if _, ok := reverse[hanzi]; !ok {
				reverse[hanzi] = id
			}
		}
	}

	return forward, reverse, rows.Err()
}

func loadSteps(ctx context.Context, db *pgxpool.Pool, lessons []lessonRow) ([]stepRow, error) {

	ids := make([]int64, len(lessons))
	for i, l := range lessons {
		ids[i] = l.ID
	}

	rows, err := db.Query(ctx, `
		SELECT lesson_id, step_type, entity_id, sort_order
		FROM lesson_steps
		WHERE lesson_id = ANY($1)
		  AND entity_id IS NOT NULL
		ORDER BY lesson_id, sort_order
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]stepRow, 0)
	for rows.Next() {
		var s stepRow
		if err := rows.Scan(&s.LessonID, &s.StepType, &s.EntityID, &s.SortOrder); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

// pickDistractors returns up to n translation strings, distinct from each
// other and from excludeID's own translation, drawn from the whole word
// pool (matching cmd/seed's forward-quiz distractor scope).
func pickDistractors(rng *rand.Rand, allWordIDs []int64, translationByID map[int64]string, excludeID int64, n int) []string {

	shuffled := append([]int64{}, allWordIDs...)
	rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	seen := map[string]bool{translationByID[excludeID]: true}
	result := make([]string, 0, n)

	for _, id := range shuffled {
		if id == excludeID {
			continue
		}
		t := translationByID[id]
		if seen[t] {
			continue
		}
		seen[t] = true
		result = append(result, t)
		if len(result) == n {
			break
		}
	}

	return result
}

// pickDistractorWords returns up to n words distinct by hanzi from the given
// id pool, excluding excludeID (matching cmd/add-lesson-variety's
// reverse-quiz distractor scope: same HSK level as the tested word).
func pickDistractorWords(rng *rand.Rand, ids []int64, byID map[int64]wordInfo, excludeID int64, n int) []wordInfo {

	shuffled := append([]int64{}, ids...)
	rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	seen := map[string]bool{byID[excludeID].Hanzi: true}
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

func createForwardQuiz(ctx context.Context, db *pgxpool.Pool, w wordInfo, distractors []string, rng *rand.Rand) (int64, error) {

	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	question := fmt.Sprintf("Как переводится %s?", w.Hanzi)

	var quizID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO quizzes (question, hsk_level, direction, hanzi, pinyin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, question, w.HSKLevel, directionWordToTranslation, w.Hanzi, w.Pinyin).Scan(&quizID); err != nil {
		return 0, err
	}

	options := append([]string{w.Translation}, distractors...)
	rng.Shuffle(len(options), func(a, b int) { options[a], options[b] = options[b], options[a] })

	for idx, optText := range options {
		if _, err := tx.Exec(ctx, `
			INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
			VALUES ($1, $2, $3, $4)
		`, quizID, optText, optText == w.Translation, idx+1); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return quizID, nil
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

func createPracticeLesson(
	ctx context.Context,
	db *pgxpool.Pool,
	courseID int64,
	lessonNumber int,
	title string,
	description string,
	steps []plannedStep,
) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var lessonID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, lesson_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, courseID, title, description, lessonNumber).Scan(&lessonID); err != nil {
		return err
	}

	for i, s := range steps {
		if _, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, $2, $3, $4)
		`, lessonID, s.StepType, s.EntityID, i+1); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
