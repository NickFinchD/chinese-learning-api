package courses

import (
	"context"
	"strconv"

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
// AdminList is independent of List: no per-user progress join, has a total
// count for pagination, and orders newest-first so new content is easy to
// find instead of buried by sort_order.
func (r *Repository) AdminList(ctx context.Context, request AdminListRequest) ([]Course, int, error) {

	page := request.Page
	if page < 1 {
		page = 1
	}

	limit := request.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	where := `WHERE 1=1`
	args := []interface{}{}
	arg := 1

	if request.Search != "" {
		where += ` AND title ILIKE $` + strconv.Itoa(arg)
		args = append(args, "%"+request.Search+"%")
		arg++
	}

	var total int

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM courses `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, title, description, hsk_level, sort_order, created_at, updated_at
		FROM courses
	` + where + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	courses := make([]Course, 0)

	for rows.Next() {

		var course Course

		if err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.HSKLevel,
			&course.SortOrder,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *Repository) Create(ctx context.Context, req CourseRequest) (*Course, error) {

	course := &Course{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO courses (title, description, hsk_level, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, hsk_level, sort_order, created_at, updated_at
	`, req.Title, req.Description, req.HSKLevel, req.SortOrder).Scan(
		&course.ID, &course.Title, &course.Description, &course.HSKLevel, &course.SortOrder,
		&course.CreatedAt, &course.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req CourseRequest) (*Course, error) {

	course := &Course{}

	err := r.db.QueryRow(ctx, `
		UPDATE courses
		SET title = $1, description = $2, hsk_level = $3, sort_order = $4, updated_at = now()
		WHERE id = $5
		RETURNING id, title, description, hsk_level, sort_order, created_at, updated_at
	`, req.Title, req.Description, req.HSKLevel, req.SortOrder, id).Scan(
		&course.ID, &course.Title, &course.Description, &course.HSKLevel, &course.SortOrder,
		&course.CreatedAt, &course.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	// lessons.course_id and lesson_steps.lesson_id both cascade, so this
	// takes the course's lessons and their steps with it.
	_, err := r.db.Exec(ctx, `DELETE FROM courses WHERE id = $1`, id)
	return err
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
