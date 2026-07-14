package courses

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context) ([]Course, error) {

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
		ORDER BY sort_order, id
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var courses []Course

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

	var lessons []LessonDTO

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
