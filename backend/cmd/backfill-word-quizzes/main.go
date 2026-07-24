// Command backfill-word-quizzes ensures every 'word' step in every lesson
// (any course) has a matching "Как переводится X?" quiz step within that
// same lesson, appending any that are missing. It exists because the
// hand-authored HSK1 lessons 1..6 (migration 000016) only ever got a single
// quiz each despite introducing up to 12 words — the programmatic
// generators (cmd/seed, cmd/seed-hsk2) guarantee full coverage themselves,
// so this is mainly a safety net for hand-authored content, but it's safe
// (and a no-op) to run against any course.
//
//	go run ./cmd/backfill-word-quizzes
//
// Purely additive and idempotent: it never removes or reorders existing
// steps, only appends missing quiz steps to the end of each lesson.
// Vocabulary quizzes are looked up by question text and reused across
// lessons instead of being duplicated.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type wordInfo struct {
	ID          int64
	Hanzi       string
	Translation string
	HSKLevel    int16
}

type lessonStepRow struct {
	LessonID  int64
	StepType  string
	EntityID  int64
	SortOrder int
	Hanzi     *string
	HSKLevel  *int16
	Question  *string
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	allWords, err := loadAllWords(ctx, db)
	if err != nil {
		log.Fatalf("failed to load words: %v", err)
	}

	wordByID := make(map[int64]wordInfo, len(allWords))
	wordByHanzi := make(map[string]wordInfo, len(allWords))
	translationByID := make(map[int64]string, len(allWords))
	wordIDsByLevel := map[int16][]int64{}

	for _, w := range allWords {
		wordByID[w.ID] = w
		wordByHanzi[w.Hanzi] = w
		translationByID[w.ID] = w.Translation
		wordIDsByLevel[w.HSKLevel] = append(wordIDsByLevel[w.HSKLevel], w.ID)
	}

	rows, err := loadLessonSteps(ctx, db)
	if err != nil {
		log.Fatalf("failed to load lesson steps: %v", err)
	}

	type lessonInfo struct {
		words        []wordInfo
		coveredWords map[int64]bool
		maxSortOrder int
	}

	lessons := map[int64]*lessonInfo{}

	for _, r := range rows {

		li, ok := lessons[r.LessonID]
		if !ok {
			li = &lessonInfo{coveredWords: map[int64]bool{}}
			lessons[r.LessonID] = li
		}

		if r.SortOrder > li.maxSortOrder {
			li.maxSortOrder = r.SortOrder
		}

		switch r.StepType {
		case "word":
			if w, ok := wordByID[r.EntityID]; ok {
				li.words = append(li.words, w)
			}
		case "quiz":
			if r.Question != nil {
				if w, ok := wordByHanziQuestion(*r.Question, wordByHanzi); ok {
					li.coveredWords[w.ID] = true
				}
			}
		}
	}

	rng := rand.New(rand.NewSource(7))
	vocabQuizCache := map[int64]int64{}
	lessonsChanged := 0
	quizzesAdded := 0

	for lessonID, li := range lessons {

		var missing []wordInfo

		for _, w := range li.words {
			if !li.coveredWords[w.ID] {
				missing = append(missing, w)
			}
		}

		if len(missing) == 0 {
			continue
		}

		sortOrder := li.maxSortOrder + 1

		tx, err := db.Begin(ctx)
		if err != nil {
			log.Fatalf("failed to begin tx for lesson %d: %v", lessonID, err)
		}

		for _, w := range missing {

			quizID, ok := vocabQuizCache[w.ID]

			if !ok {

				question := fmt.Sprintf("Как переводится %s?", w.Hanzi)

				existingID, found, err := findQuizByQuestion(ctx, tx, question)
				if err != nil {
					tx.Rollback(ctx)
					log.Fatalf("failed to look up quiz for %q: %v", w.Hanzi, err)
				}

				if found {
					quizID = existingID
				} else {

					distractors := pickDistractors(rng, wordIDsByLevel[w.HSKLevel], translationByID, w.ID, 3)

					if err := tx.QueryRow(ctx, `
						INSERT INTO quizzes (question, hsk_level) VALUES ($1, $2) RETURNING id
					`, question, w.HSKLevel).Scan(&quizID); err != nil {
						tx.Rollback(ctx)
						log.Fatalf("failed to create quiz for %q: %v", w.Hanzi, err)
					}

					options := append([]string{w.Translation}, distractors...)

					rng.Shuffle(len(options), func(a, b int) {
						options[a], options[b] = options[b], options[a]
					})

					for idx, optText := range options {

						_, err := tx.Exec(ctx, `
							INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
							VALUES ($1, $2, $3, $4)
						`, quizID, optText, optText == w.Translation, idx+1)
						if err != nil {
							tx.Rollback(ctx)
							log.Fatalf("failed to create quiz options for %q: %v", w.Hanzi, err)
						}
					}
				}

				vocabQuizCache[w.ID] = quizID
			}

			if _, err := tx.Exec(ctx, `
				INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
				VALUES ($1, 'quiz', $2, $3)
			`, lessonID, quizID, sortOrder); err != nil {
				tx.Rollback(ctx)
				log.Fatalf("failed to insert quiz step for lesson %d: %v", lessonID, err)
			}

			sortOrder++
			quizzesAdded++
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("failed to commit lesson %d: %v", lessonID, err)
		}

		lessonsChanged++
		fmt.Printf("Lesson %d: added %d missing quiz(zes)\n", lessonID, len(missing))
	}

	fmt.Printf("Done. %d lesson(s) updated, %d quiz step(s) added.\n", lessonsChanged, quizzesAdded)
}

// wordByHanziQuestion reverses "Как переводится X?" back to the word it
// covers. Auto-generated quizzes always follow this exact format.
func wordByHanziQuestion(question string, wordByHanzi map[string]wordInfo) (wordInfo, bool) {
	const prefix = "Как переводится "
	const suffix = "?"

	if len(question) < len(prefix)+len(suffix) {
		return wordInfo{}, false
	}

	if question[:len(prefix)] != prefix || question[len(question)-len(suffix):] != suffix {
		return wordInfo{}, false
	}

	hanzi := question[len(prefix) : len(question)-len(suffix)]

	w, ok := wordByHanzi[hanzi]

	return w, ok
}

func loadAllWords(ctx context.Context, db *pgxpool.Pool) ([]wordInfo, error) {

	rows, err := db.Query(ctx, `SELECT id, hanzi, translation, hsk_level FROM words ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]wordInfo, 0)

	for rows.Next() {

		var w wordInfo

		if err := rows.Scan(&w.ID, &w.Hanzi, &w.Translation, &w.HSKLevel); err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	return words, rows.Err()
}

func loadLessonSteps(ctx context.Context, db *pgxpool.Pool) ([]lessonStepRow, error) {

	rows, err := db.Query(ctx, `
		SELECT
			ls.lesson_id,
			ls.step_type,
			ls.entity_id,
			ls.sort_order,
			w.hanzi,
			w.hsk_level,
			q.question
		FROM lesson_steps ls
		LEFT JOIN words w ON w.id = ls.entity_id AND ls.step_type = 'word'
		LEFT JOIN quizzes q ON q.id = ls.entity_id AND ls.step_type = 'quiz'
		ORDER BY ls.lesson_id, ls.sort_order
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]lessonStepRow, 0)

	for rows.Next() {

		var r lessonStepRow

		if err := rows.Scan(&r.LessonID, &r.StepType, &r.EntityID, &r.SortOrder, &r.Hanzi, &r.HSKLevel, &r.Question); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, rows.Err()
}

func pickDistractors(rng *rand.Rand, allWordIDs []int64, translationByID map[int64]string, excludeID int64, n int) []string {

	shuffled := append([]int64{}, allWordIDs...)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

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

func findQuizByQuestion(ctx context.Context, tx pgx.Tx, question string) (int64, bool, error) {

	var id int64

	err := tx.QueryRow(ctx, `SELECT id FROM quizzes WHERE question = $1`, question).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}

	return id, true, nil
}
