// Command seed-hsk3-tests creates one standalone "Как переводится X?" vocab
// quiz (hsk_level = 3) per HSK3 word, for words that don't already have one.
// Unlike cmd/seed and cmd/backfill-word-quizzes, these quizzes aren't
// attached to any lesson_steps — HSK3 has no course/lessons yet, and the
// grammar-test page (GrammarTestPage.vue) reads quizzes directly by
// hsk_level, so a lesson isn't required for them to be usable.
//
//	go run ./cmd/seed-hsk3-tests
//
// Idempotent: quizzes are looked up by question text and reused, so re-runs
// only fill in gaps left by newly added HSK3 words.
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
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	words, err := loadWords(ctx, db)
	if err != nil {
		log.Fatalf("failed to load HSK3 words: %v", err)
	}

	translationByID := make(map[int64]string, len(words))
	allWordIDs := make([]int64, 0, len(words))

	for _, w := range words {
		translationByID[w.ID] = w.Translation
		allWordIDs = append(allWordIDs, w.ID)
	}

	rng := rand.New(rand.NewSource(3))
	created := 0
	skipped := 0

	for _, w := range words {

		question := fmt.Sprintf("Как переводится %s?", w.Hanzi)

		tx, err := db.Begin(ctx)
		if err != nil {
			log.Fatalf("failed to begin tx: %v", err)
		}

		_, found, err := findQuizByQuestion(ctx, tx, question)
		if err != nil {
			tx.Rollback(ctx)
			log.Fatalf("failed to look up quiz for %q: %v", w.Hanzi, err)
		}

		if found {
			tx.Rollback(ctx)
			skipped++
			continue
		}

		var quizID int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO quizzes (question, hsk_level) VALUES ($1, 3) RETURNING id
		`, question).Scan(&quizID); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("failed to create quiz for %q: %v", w.Hanzi, err)
		}

		distractors := pickDistractors(rng, allWordIDs, translationByID, w.ID, 3)
		options := append([]string{w.Translation}, distractors...)

		rng.Shuffle(len(options), func(a, b int) {
			options[a], options[b] = options[b], options[a]
		})

		for idx, optText := range options {
			if _, err := tx.Exec(ctx, `
				INSERT INTO quiz_options (quiz_id, option_text, is_correct, sort_order)
				VALUES ($1, $2, $3, $4)
			`, quizID, optText, optText == w.Translation, idx+1); err != nil {
				tx.Rollback(ctx)
				log.Fatalf("failed to create quiz options for %q: %v", w.Hanzi, err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("failed to commit quiz for %q: %v", w.Hanzi, err)
		}

		created++
	}

	fmt.Printf("Done. %d quiz(zes) created, %d already existed.\n", created, skipped)
}

func loadWords(ctx context.Context, db *pgxpool.Pool) ([]wordInfo, error) {

	rows, err := db.Query(ctx, `SELECT id, hanzi, translation FROM words WHERE hsk_level = 3 ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]wordInfo, 0)

	for rows.Next() {

		var w wordInfo

		if err := rows.Scan(&w.ID, &w.Hanzi, &w.Translation); err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	return words, rows.Err()
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
