package courses

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

func (r *Repository) List(ctx context.Context, userID int64) ([]Course, error) {

	query := `
		SELECT
			c.id,
			c.title,
			c.description,
			c.hsk_level,
			c.sort_order,
			c.created_at,
			c.updated_at,
			COALESCE(ucp.progress_percent, 0)
		FROM courses c
		LEFT JOIN user_course_progress ucp
			ON ucp.course_id = c.id AND ucp.user_id = $1
		ORDER BY c.sort_order, c.id
	`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	courses := make([]Course, 0)

	for rows.Next() {

		var course Course

		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.HSKLevel,
			&course.SortOrder,
			&course.CreatedAt,
			&course.UpdatedAt,
			&course.ProgressPercent,
		)

		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	return courses, rows.Err()
}
func (r *Repository) GetByID(ctx context.Context, id int64) (*Course, error) {

	query := `
		SELECT
			id,
			title,
			description,
			hsk_level,
			sort_order,
			created_at,
			updated_at
		FROM courses
		WHERE id = $1
	`

	course := &Course{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.HSKLevel,
		&course.SortOrder,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return course, nil
}
func (r *Repository) GetLessons(ctx context.Context, courseID int64) ([]LessonDTO, error) {

	query := `
		SELECT
			id,
			title,
			lesson_number
		FROM lessons
		WHERE course_id = $1
		ORDER BY lesson_number
	`

	rows, err := r.db.Query(ctx, query, courseID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	lessons := make([]LessonDTO, 0)

	for rows.Next() {

		var lesson LessonDTO

		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.LessonNumber,
		)

		if err != nil {
			return nil, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, rows.Err()
}
