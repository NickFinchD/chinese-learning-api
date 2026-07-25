// Command fix-sentence-builder-vocab corrects lesson sentence_builder steps
// that use words the learner hasn't been taught yet in that course. The
// original placement (cmd/add-lesson-variety) matched sentence_exercises
// only by HSK level, but most of that pool's sentences use vocabulary well
// beyond the basic word list regardless of level — only 23 of 165 existing
// exercises use exclusively catalogued words at all, so lesson position
// never factored in.
//
// This recomputes, for every eligible lesson (lesson_number >= 3, excluding
// lesson 995 which has its own hand-authored sentence_builder steps predating
// this whole feature), the learner's cumulative vocabulary up to and
// including that lesson — for HSK2, that's the full HSK1 word list plus
// HSK2 words up to the lesson, since the course continues from HSK1 — and
// only assigns an exercise whose every chunk is a word already taught by
// then. If no such exercise exists, the lesson is left without one rather
// than forcing inappropriate content in.
//
//	go run ./cmd/fix-sentence-builder-vocab
//
// Safe to rerun: recomputes the same way every time.
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

const excludedLessonID = 995 // pre-existing hand-authored review lesson

type lessonRow struct {
	ID           int64
	CourseID     int64
	LessonNumber int
}

type exercise struct {
	ID     int64
	Chunks []string
}

type stepRow struct {
	ID        int64
	LessonID  int64
	StepType  string
	EntityID  int64
	SortOrder int
}

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	seeded, err := seedNewSentences(ctx, db)
	if err != nil {
		log.Fatalf("failed to seed new sentences: %v", err)
	}
	fmt.Printf("Seeded %d new sentence_exercises (skipped ones already present).\n", seeded)

	lessons, err := loadLessons(ctx, db)
	if err != nil {
		log.Fatalf("failed to load lessons: %v", err)
	}

	wordSteps, err := loadWordSteps(ctx, db)
	if err != nil {
		log.Fatalf("failed to load word steps: %v", err)
	}

	exercises, err := loadExercises(ctx, db)
	if err != nil {
		log.Fatalf("failed to load sentence exercises: %v", err)
	}

	steps, err := loadSteps(ctx, db)
	if err != nil {
		log.Fatalf("failed to load lesson steps: %v", err)
	}

	byLesson := map[int64][]stepRow{}
	for _, s := range steps {
		byLesson[s.LessonID] = append(byLesson[s.LessonID], s)
	}

	// HSK1 course id is assumed to be the one whose lessons feed HSK2's
	// baseline vocabulary; resolved dynamically as "the other course".
	lessonsByCourse := map[int64][]lessonRow{}
	for _, l := range lessons {
		lessonsByCourse[l.CourseID] = append(lessonsByCourse[l.CourseID], l)
	}
	for courseID := range lessonsByCourse {
		sort.Slice(lessonsByCourse[courseID], func(i, j int) bool {
			return lessonsByCourse[courseID][i].LessonNumber < lessonsByCourse[courseID][j].LessonNumber
		})
	}

	courseIDs := make([]int64, 0, len(lessonsByCourse))
	for id := range lessonsByCourse {
		courseIDs = append(courseIDs, id)
	}
	sort.Slice(courseIDs, func(i, j int) bool { return courseIDs[i] < courseIDs[j] })

	if len(courseIDs) != 2 {
		log.Fatalf("expected exactly 2 courses, found %d", len(courseIDs))
	}

	baseCourseID, advancedCourseID := courseIDs[0], courseIDs[1]

	// Full cumulative vocab of the base (HSK1) course, used as the starting
	// point for the advanced (HSK2) course's own cumulative vocab.
	baseFullVocab := map[string]bool{}
	for _, l := range lessonsByCourse[baseCourseID] {
		for _, hanzi := range wordSteps[l.ID] {
			baseFullVocab[hanzi] = true
		}
	}

	rng := rand.New(rand.NewSource(11))

	unchanged, replaced, removed, addedNone := 0, 0, 0, 0

	for _, courseID := range courseIDs {

		cumulative := map[string]bool{}
		if courseID == advancedCourseID {
			for hanzi := range baseFullVocab {
				cumulative[hanzi] = true
			}
		}

		for _, lesson := range lessonsByCourse[courseID] {

			for _, hanzi := range wordSteps[lesson.ID] {
				cumulative[hanzi] = true
			}

			if lesson.ID == excludedLessonID || lesson.LessonNumber < 3 {
				continue
			}

			original := byLesson[lesson.ID]
			sort.Slice(original, func(i, j int) bool { return original[i].SortOrder < original[j].SortOrder })

			var current *stepRow
			for i := range original {
				if original[i].StepType == "sentence_builder" {
					current = &original[i]
					break
				}
			}

			if current == nil {
				continue
			}

			currentValid := exercises[current.EntityID] != nil && allChunksKnown(exercises[current.EntityID].Chunks, cumulative)

			if currentValid {
				unchanged++
				continue
			}

			candidates := make([]*exercise, 0)
			for _, ex := range exercises {
				if allChunksKnown(ex.Chunks, cumulative) {
					candidates = append(candidates, ex)
				}
			}

			if len(candidates) == 0 {
				if err := deleteStep(ctx, db, lesson.ID, current.ID); err != nil {
					log.Fatalf("failed to remove invalid sentence_builder from lesson %d: %v", lesson.ID, err)
				}
				removed++
				addedNone++
				continue
			}

			sort.Slice(candidates, func(i, j int) bool { return candidates[i].ID < candidates[j].ID })
			pick := candidates[rng.Intn(len(candidates))]

			if err := updateStepEntity(ctx, db, current.ID, pick.ID); err != nil {
				log.Fatalf("failed to update sentence_builder for lesson %d: %v", lesson.ID, err)
			}
			replaced++
		}
	}

	fmt.Printf(
		"Done. %d unchanged (already valid), %d replaced, %d removed (no valid replacement).\n",
		unchanged, replaced, removed,
	)
}

// seedNewSentences inserts newSentences that aren't already present
// (matched by translation), so the command is safe to rerun against a
// freshly migrated database.
func seedNewSentences(ctx context.Context, db *pgxpool.Pool) (int, error) {

	existing := map[string]bool{}

	rows, err := db.Query(ctx, `SELECT translation FROM sentence_exercises`)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			rows.Close()
			return 0, err
		}
		existing[t] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	seeded := 0
	for _, s := range newSentences {

		if existing[s.Translation] {
			continue
		}

		if _, err := db.Exec(ctx, `
			INSERT INTO sentence_exercises (translation, chunks, pinyin, hsk_level)
			VALUES ($1, $2, $3, $4)
		`, s.Translation, s.Chunks, s.Pinyin, s.HSKLevel); err != nil {
			return seeded, err
		}
		seeded++
	}

	return seeded, nil
}

func allChunksKnown(chunks []string, known map[string]bool) bool {
	for _, c := range chunks {
		if !known[c] {
			return false
		}
	}
	return true
}

func loadLessons(ctx context.Context, db *pgxpool.Pool) ([]lessonRow, error) {
	rows, err := db.Query(ctx, `SELECT id, course_id, lesson_number FROM lessons`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]lessonRow, 0)
	for rows.Next() {
		var l lessonRow
		if err := rows.Scan(&l.ID, &l.CourseID, &l.LessonNumber); err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, rows.Err()
}

// loadWordSteps returns, per lesson, the hanzi of every word introduced via
// a 'word' step in that lesson.
func loadWordSteps(ctx context.Context, db *pgxpool.Pool) (map[int64][]string, error) {
	rows, err := db.Query(ctx, `
		SELECT ls.lesson_id, w.hanzi
		FROM lesson_steps ls
		JOIN words w ON w.id = ls.entity_id
		WHERE ls.step_type = 'word'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]string)
	for rows.Next() {
		var lessonID int64
		var hanzi string
		if err := rows.Scan(&lessonID, &hanzi); err != nil {
			return nil, err
		}
		result[lessonID] = append(result[lessonID], hanzi)
	}
	return result, rows.Err()
}

func loadExercises(ctx context.Context, db *pgxpool.Pool) (map[int64]*exercise, error) {
	rows, err := db.Query(ctx, `SELECT id, chunks FROM sentence_exercises`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]*exercise)
	for rows.Next() {
		var ex exercise
		if err := rows.Scan(&ex.ID, &ex.Chunks); err != nil {
			return nil, err
		}
		result[ex.ID] = &ex
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

func updateStepEntity(ctx context.Context, db *pgxpool.Pool, stepID, newEntityID int64) error {
	_, err := db.Exec(ctx, `UPDATE lesson_steps SET entity_id = $1 WHERE id = $2`, newEntityID, stepID)
	return err
}

func deleteStep(ctx context.Context, db *pgxpool.Pool, lessonID, stepID int64) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM lesson_steps WHERE id = $1`, stepID); err != nil {
		return err
	}

	rows, err := tx.Query(ctx, `
		SELECT id FROM lesson_steps WHERE lesson_id = $1 ORDER BY sort_order
	`, lessonID)
	if err != nil {
		return err
	}

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE lesson_steps SET sort_order = $1 WHERE id = $2`, i+1, id); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
