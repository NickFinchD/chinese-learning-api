package lessons

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
