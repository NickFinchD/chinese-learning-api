package progress

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) StartLesson(
	ctx context.Context,
	userID int64,
	lessonID int64,
) error {

	query := `
		INSERT INTO user_lesson_progress (
			user_id,
			lesson_id,
			status,
			current_step,
			started_at
		)
		VALUES (
			$1,
			$2,
			'in_progress',
			1,
			NOW()
		)
		ON CONFLICT (user_id, lesson_id)
		DO UPDATE
		SET
			status = 'in_progress',
			current_step = 1,
			started_at = NOW(),
			updated_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, lessonID)

	return err
}
func (r *Repository) GetProgress(
	ctx context.Context,
	userID int64,
	lessonID int64,
) (*UserLessonProgress, error) {

	query := `
		SELECT
			id,
			user_id,
			lesson_id,
			status,
			current_step,
			score,
			started_at,
			completed_at,
			created_at,
			updated_at
		FROM user_lesson_progress
		WHERE user_id = $1
		AND lesson_id = $2
	`

	progress := &UserLessonProgress{}

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		lessonID,
	).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.LessonID,
		&progress.Status,
		&progress.CurrentStep,
		&progress.Score,
		&progress.StartedAt,
		&progress.CompletedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return progress, nil
}
func (r *Repository) UpdateStep(
	ctx context.Context,
	userID int64,
	lessonID int64,
	currentStep int,
) error {

	query := `
		UPDATE user_lesson_progress
		SET
			current_step = $3,
			updated_at = NOW()
		WHERE
			user_id = $1
			AND lesson_id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		lessonID,
		currentStep,
	)

	return err
}
func (r *Repository) CompleteLesson(
	ctx context.Context,
	userID int64,
	lessonID int64,
	score int,
) error {

	query := `
		UPDATE user_lesson_progress
		SET
			status = 'completed',
			score = $3,
			completed_at = NOW(),
			updated_at = NOW()
		WHERE
			user_id = $1
			AND lesson_id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		lessonID,
		score,
	)

	return err
}
func (r *Repository) UpdateCourseProgress(
	ctx context.Context,
	userID int64,
	lessonID int64,
) error {

	query := `
		WITH lesson_data AS (
			SELECT course_id
			FROM lessons
			WHERE id = $2
		),
		total AS (
			SELECT COUNT(*) AS total_lessons
			FROM lessons
			WHERE course_id = (SELECT course_id FROM lesson_data)
		),
		completed AS (
			SELECT COUNT(*) AS completed_lessons
			FROM user_lesson_progress ulp
			JOIN lessons l ON l.id = ulp.lesson_id
			WHERE ulp.user_id = $1
				AND ulp.status = 'completed'
				AND l.course_id = (SELECT course_id FROM lesson_data)
		)

		INSERT INTO user_course_progress (
			user_id,
			course_id,
			completed_lessons,
			total_lessons,
			progress_percent,
			updated_at
		)

		SELECT
			$1,
			(SELECT course_id FROM lesson_data),
			COALESCE((SELECT completed_lessons FROM completed), 0),
			(SELECT total_lessons FROM total),
			ROUND(
				COALESCE((SELECT completed_lessons FROM completed), 0) * 100.0 /
				(SELECT total_lessons FROM total)
			),
			NOW()

		ON CONFLICT (user_id, course_id)
		DO UPDATE SET
			completed_lessons = EXCLUDED.completed_lessons,
			total_lessons = EXCLUDED.total_lessons,
			progress_percent = EXCLUDED.progress_percent,
			updated_at = NOW();
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		lessonID,
	)

	return err
}
