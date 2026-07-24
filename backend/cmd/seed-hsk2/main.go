// Command seed-hsk2 generates all lessons for the "HSK 2" course from the
// 319 HSK2 words already seeded in migration 000020 (there is no
// hand-authored content for this course, unlike HSK1 lessons 1..6).
//
// Words are first grouped by part of speech (grammar words like numerals/
// pronouns/particles first, then adverbs/verbs/adjectives, then nouns last —
// see categoryOrder), then introduced in small batches within each group, so
// each lesson gets a real, descriptive title ("Существительные 3",
// "Глаголы 1", ...) instead of a bare "Часть N". Every few batches a
// dedicated review-only lesson ("Повторение N") is inserted with no new
// words, just extra word/quiz practice drawn from everything introduced so
// far — this is what stretches the course length without diluting content.
// Vocabulary quizzes are auto-generated ("Как переводится X?"). There is no
// grammar-quiz bank for HSK2 (that content is hand-authored per level and
// out of scope here).
//
// It is meant to be run once against a database that already has migration
// 000023 applied:
//
//	go run ./cmd/seed-hsk2
//
// It is destructive-but-idempotent with respect to the HSK2 course: it
// deletes all of that course's lessons (and their steps/progress, via ON
// DELETE CASCADE) before regenerating them, so it can be re-run freely.
// Vocabulary quizzes are looked up by question text and reused instead of
// being duplicated.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/database"
	"github.com/jackc/pgx/v5"
)

const (
	hskLevel               = 2
	newWordsPerLesson      = 3
	wordStepsPerLesson     = 12
	bonusReviewQuizzes     = 2 // extra quizzes from review words, on top of the guaranteed one-per-new-word
	batchesPerReviewLesson = 4
	reviewLessonWordCount  = 10
)

// categoryOrder controls both the grouping (grammar words first, content
// words last) and, indirectly, lesson titles: consecutive batches sharing a
// dominant part of speech get numbered ("Глаголы 1", "Глаголы 2", ...).
var categoryOrder = []string{
	"numeral", "pronoun", "interrogative", "particle",
	"measure_word", "preposition", "conjunction", "adverb",
	"verb", "adjective", "noun", "phrase",
}

var categoryLabel = map[string]string{
	"numeral":       "Числительные",
	"pronoun":       "Местоимения",
	"interrogative": "Вопросительные слова",
	"particle":      "Частицы",
	"measure_word":  "Счётные слова",
	"preposition":   "Предлоги",
	"conjunction":   "Союзы",
	"adverb":        "Наречия",
	"verb":          "Глаголы",
	"adjective":     "Прилагательные",
	"noun":          "Существительные",
	"phrase":        "Устойчивые выражения",
}

func categoryRank(pos string) int {
	for i, c := range categoryOrder {
		if c == pos {
			return i
		}
	}
	return len(categoryOrder)
}

func categoryTitle(pos string) string {
	if label, ok := categoryLabel[pos]; ok {
		return label
	}
	return "Смешанная лексика"
}

type wordInfo struct {
	ID           int64
	Hanzi        string
	Translation  string
	PartOfSpeech string
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	var courseID int64

	if err := db.QueryRow(ctx, `SELECT id FROM courses WHERE title = 'HSK 2'`).Scan(&courseID); err != nil {
		log.Fatalf("failed to find HSK 2 course: %v", err)
	}

	if _, err := db.Exec(ctx, `DELETE FROM lessons WHERE course_id = $1`, courseID); err != nil {
		log.Fatalf("failed to clear previous HSK2 lessons: %v", err)
	}

	allWords, err := loadWords(ctx, db)
	if err != nil {
		log.Fatalf("failed to load HSK2 words: %v", err)
	}

	if len(allWords) == 0 {
		log.Fatal("no HSK2 words found — run migrations up to 000020 first")
	}

	sort.SliceStable(allWords, func(i, j int) bool {
		return categoryRank(allWords[i].PartOfSpeech) < categoryRank(allWords[j].PartOfSpeech)
	})

	translationByID := make(map[int64]string, len(allWords))
	allWordIDs := make([]int64, 0, len(allWords))

	for _, w := range allWords {
		translationByID[w.ID] = w.Translation
		allWordIDs = append(allWordIDs, w.ID)
	}

	rng := rand.New(rand.NewSource(2))

	vocabQuizCache := map[int64]int64{}
	introduced := make([]wordInfo, 0, len(allWords))
	lessonNum := 1
	batchCount := 0
	reviewCount := 0
	lastCategory := ""
	categoryPart := 0

	for start := 0; start < len(allWords); start += newWordsPerLesson {

		end := start + newWordsPerLesson
		if end > len(allWords) {
			end = len(allWords)
		}

		newWords := allWords[start:end]

		pool := append([]wordInfo{}, introduced...)

		rng.Shuffle(len(pool), func(i, j int) {
			pool[i], pool[j] = pool[j], pool[i]
		})

		fillCount := wordStepsPerLesson - len(newWords)
		if fillCount > len(pool) {
			fillCount = len(pool)
		}
		if fillCount < 0 {
			fillCount = 0
		}

		lessonWords := append(append([]wordInfo{}, newWords...), pool[:fillCount]...)

		category := dominantCategory(newWords)

		if category == lastCategory {
			categoryPart++
		} else {
			categoryPart = 1
			lastCategory = category
		}

		title := categoryTitle(category)
		if categoryPart > 1 {
			title = fmt.Sprintf("%s %d", title, categoryPart)
		}

		hanziList := ""
		for i, w := range newWords {
			if i > 0 {
				hanziList += ", "
			}
			hanziList += w.Hanzi
		}

		description := fmt.Sprintf("Новые слова: %s", hanziList)

		if err := createLesson(ctx, db, courseID, lessonNum, title, description, lessonWords,
			newWords, bonusReviewQuizzes, vocabQuizCache, translationByID, allWordIDs, rng); err != nil {
			log.Fatalf("failed to create lesson %d: %v", lessonNum, err)
		}

		fmt.Printf("Created lesson %d: %s (%d new, %d review)\n", lessonNum, title, len(newWords), fillCount)

		introduced = append(introduced, newWords...)
		lessonNum++
		batchCount++

		if batchCount%batchesPerReviewLesson == 0 {

			reviewCount++

			reviewPool := append([]wordInfo{}, introduced...)

			rng.Shuffle(len(reviewPool), func(i, j int) {
				reviewPool[i], reviewPool[j] = reviewPool[j], reviewPool[i]
			})

			n := reviewLessonWordCount
			if n > len(reviewPool) {
				n = len(reviewPool)
			}

			reviewWords := reviewPool[:n]
			reviewTitle := fmt.Sprintf("Повторение %d", reviewCount)
			reviewDescription := "Повторение изученных слов: больше практики без новой лексики"

			if err := createLesson(ctx, db, courseID, lessonNum, reviewTitle, reviewDescription, reviewWords,
				reviewWords, 0, vocabQuizCache, translationByID, allWordIDs, rng); err != nil {
				log.Fatalf("failed to create review lesson %d: %v", lessonNum, err)
			}

			fmt.Printf("Created lesson %d: %s (%d review words)\n", lessonNum, reviewTitle, n)

			lessonNum++
		}
	}

	fmt.Printf("Done. %d lessons total.\n", lessonNum-1)
}

func dominantCategory(words []wordInfo) string {

	counts := map[string]int{}
	order := make([]string, 0, len(words))

	for _, w := range words {
		if counts[w.PartOfSpeech] == 0 {
			order = append(order, w.PartOfSpeech)
		}
		counts[w.PartOfSpeech]++
	}

	best := ""
	bestCount := -1

	for _, pos := range order {
		if counts[pos] > bestCount {
			best = pos
			bestCount = counts[pos]
		}
	}

	return best
}

func loadWords(ctx context.Context, db pgxPool) ([]wordInfo, error) {

	rows, err := db.Query(ctx, `
		SELECT id, hanzi, translation, COALESCE(part_of_speech, '') FROM words
		WHERE hsk_level = $1
		ORDER BY id
	`, hskLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]wordInfo, 0)

	for rows.Next() {

		var w wordInfo

		if err := rows.Scan(&w.ID, &w.Hanzi, &w.Translation, &w.PartOfSpeech); err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	return words, rows.Err()
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

// pgxPool is the subset of *pgxpool.Pool used by the loader helpers above.
type pgxPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
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

func createLesson(
	ctx context.Context,
	db interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	},
	courseID int64,
	lessonNum int,
	title string,
	description string,
	lessonWords []wordInfo,
	newWords []wordInfo, // every one of these gets a quiz — guaranteed practice for what's actually taught this lesson
	bonusQuizCount int, // extra quizzes drawn from the remaining (review) words
	vocabQuizCache map[int64]int64,
	translationByID map[int64]string,
	allWordIDs []int64,
	rng *rand.Rand,
) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var lessonID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, lesson_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, courseID, title, description, lessonNum).Scan(&lessonID)
	if err != nil {
		return err
	}

	sortOrder := 1

	for _, w := range lessonWords {

		_, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'word', $2, $3)
		`, lessonID, w.ID, sortOrder)
		if err != nil {
			return err
		}

		sortOrder++
	}

	newWordIDs := make(map[int64]bool, len(newWords))
	for _, w := range newWords {
		newWordIDs[w.ID] = true
	}

	reviewOnly := make([]wordInfo, 0, len(lessonWords))
	for _, w := range lessonWords {
		if !newWordIDs[w.ID] {
			reviewOnly = append(reviewOnly, w)
		}
	}

	rng.Shuffle(len(reviewOnly), func(i, j int) {
		reviewOnly[i], reviewOnly[j] = reviewOnly[j], reviewOnly[i]
	})

	if bonusQuizCount > len(reviewOnly) {
		bonusQuizCount = len(reviewOnly)
	}

	// Every new word is guaranteed a quiz; bonus review quizzes are appended
	// on top so practice never falls short of what was just taught.
	quizCandidates := append(append([]wordInfo{}, newWords...), reviewOnly[:bonusQuizCount]...)

	for i := 0; i < len(quizCandidates); i++ {

		w := quizCandidates[i]

		quizID, ok := vocabQuizCache[w.ID]

		if !ok {

			question := fmt.Sprintf("Как переводится %s?", w.Hanzi)

			existingID, found, err := findQuizByQuestion(ctx, tx, question)
			if err != nil {
				return err
			}

			if found {
				quizID = existingID
			} else {

				distractors := pickDistractors(rng, allWordIDs, translationByID, w.ID, 3)

				if err := tx.QueryRow(ctx, `
					INSERT INTO quizzes (question, hsk_level) VALUES ($1, $2) RETURNING id
				`, question, hskLevel).Scan(&quizID); err != nil {
					return err
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
						return err
					}
				}
			}

			vocabQuizCache[w.ID] = quizID
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
			VALUES ($1, 'quiz', $2, $3)
		`, lessonID, quizID, sortOrder)
		if err != nil {
			return err
		}

		sortOrder++
	}

	return tx.Commit(ctx)
}
