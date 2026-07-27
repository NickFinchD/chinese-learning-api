package lessons

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminListByCourse returns every lesson in a course (ordered by lesson
// number, unpaginated — course lesson counts are small) with a step count
// for the course's lesson list screen.
func (r *Repository) AdminListByCourse(ctx context.Context, courseID int64) ([]AdminLessonListItem, error) {

	query := `
		SELECT
			l.id, l.course_id, l.title, l.description, l.lesson_number,
			(SELECT COUNT(*) FROM lesson_steps ls WHERE ls.lesson_id = l.id)
		FROM lessons l
		WHERE l.course_id = $1
		ORDER BY l.lesson_number
	`

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]AdminLessonListItem, 0)

	for rows.Next() {

		var item AdminLessonListItem

		if err := rows.Scan(&item.ID, &item.CourseID, &item.Title, &item.Description, &item.LessonNumber, &item.StepCount); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

func (r *Repository) Create(ctx context.Context, req LessonRequest) (*Lesson, error) {

	lesson := &Lesson{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, lesson_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id, course_id, title, description, lesson_number, created_at, updated_at
	`, req.CourseID, req.Title, req.Description, req.LessonNumber).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Description, &lesson.LessonNumber,
		&lesson.CreatedAt, &lesson.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lesson, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req LessonRequest) (*Lesson, error) {

	lesson := &Lesson{}

	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		SET course_id = $1, title = $2, description = $3, lesson_number = $4, updated_at = now()
		WHERE id = $5
		RETURNING id, course_id, title, description, lesson_number, created_at, updated_at
	`, req.CourseID, req.Title, req.Description, req.LessonNumber, id).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Description, &lesson.LessonNumber,
		&lesson.CreatedAt, &lesson.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lesson, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	// lesson_steps.lesson_id cascades, so its steps go with it.
	_, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1`, id)
	return err
}

// AddStep appends a step to the end of the lesson (sort_order = current
// max + 1), so adding steps through the admin UI never needs a manual
// order value.
func (r *Repository) AddStep(ctx context.Context, lessonID int64, stepType string, entityID int64) (*LessonStep, error) {

	var sortOrder int

	if err := r.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(sort_order), 0) + 1 FROM lesson_steps WHERE lesson_id = $1
	`, lessonID).Scan(&sortOrder); err != nil {
		return nil, err
	}

	step := &LessonStep{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO lesson_steps (lesson_id, step_type, entity_id, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, lesson_id, step_type, entity_id, sort_order, created_at
	`, lessonID, stepType, entityID, sortOrder).Scan(
		&step.ID, &step.LessonID, &step.StepType, &step.EntityID, &step.SortOrder, &step.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return step, nil
}

func (r *Repository) RemoveStep(ctx context.Context, stepID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM lesson_steps WHERE id = $1`, stepID)
	return err
}

// ReorderSteps assigns sort_order from each id's position in stepIDs. The
// WHERE clause scopes every update to lessonID as well as the id, so a step
// id from a different lesson (client bug) is silently ignored rather than
// corrupting another lesson's order.
func (r *Repository) ReorderSteps(ctx context.Context, lessonID int64, stepIDs []int64) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, id := range stepIDs {
		if _, err := tx.Exec(ctx, `
			UPDATE lesson_steps SET sort_order = $1 WHERE id = $2 AND lesson_id = $3
		`, i+1, id, lessonID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Lesson, error) {

	query := `
		SELECT
			id,
			course_id,
			title,
			description,
			lesson_number,
			created_at,
			updated_at
		FROM lessons
		WHERE id = $1
	`

	lesson := &Lesson{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Title,
		&lesson.Description,
		&lesson.LessonNumber,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lesson, nil
}

func (r *Repository) GetSteps(ctx context.Context, lessonID int64) ([]LessonStep, error) {

	query := `
		SELECT
			id,
			lesson_id,
			step_type,
			entity_id,
			sort_order,
			created_at
		FROM lesson_steps
		WHERE lesson_id = $1
		ORDER BY sort_order
	`

	rows, err := r.db.Query(ctx, query, lessonID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	steps := make([]LessonStep, 0)

	for rows.Next() {

		var step LessonStep

		err := rows.Scan(
			&step.ID,
			&step.LessonID,
			&step.StepType,
			&step.EntityID,
			&step.SortOrder,
			&step.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		steps = append(steps, step)
	}

	return steps, rows.Err()
}
