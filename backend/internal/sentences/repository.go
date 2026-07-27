package sentences

import (
	"context"
	"strconv"
	"strings"

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

func (r *Repository) List(ctx context.Context, hskLevel int16) ([]Exercise, error) {

	query := `
		SELECT id, translation, chunks, pinyin, hsk_level, created_at, updated_at
		FROM sentence_exercises
	`

	args := []interface{}{}

	if hskLevel > 0 {
		query += ` WHERE hsk_level = $1`
		args = append(args, hskLevel)
	}

	query += ` ORDER BY hsk_level, id`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Exercise, 0)

	for rows.Next() {

		var e Exercise

		if err := rows.Scan(
			&e.ID,
			&e.Translation,
			&e.Chunks,
			&e.Pinyin,
			&e.HSKLevel,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, rows.Err()
}

// AdminList is independent of List: it paginates, supports search/hsk
// filtering, and orders newest-first so new content is easy to find.
func (r *Repository) AdminList(ctx context.Context, request AdminListRequest) ([]Exercise, int, error) {

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
		where += `
			AND (
				translation ILIKE $` + strconv.Itoa(arg) + `
				OR pinyin ILIKE $` + strconv.Itoa(arg) + `
			)
		`

		args = append(args, "%"+request.Search+"%")
		arg++
	}

	if request.HSK > 0 {
		where += ` AND hsk_level = $` + strconv.Itoa(arg)

		args = append(args, request.HSK)
		arg++
	}

	var total int

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM sentence_exercises `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, translation, chunks, pinyin, hsk_level, created_at, updated_at
		FROM sentence_exercises
	` + where + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	result := make([]Exercise, 0)

	for rows.Next() {

		var e Exercise

		if err := rows.Scan(
			&e.ID,
			&e.Translation,
			&e.Chunks,
			&e.Pinyin,
			&e.HSKLevel,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		result = append(result, e)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *Repository) Create(ctx context.Context, req AdminExerciseRequest) (*Exercise, error) {

	e := &Exercise{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO sentence_exercises (translation, chunks, pinyin, hsk_level)
		VALUES ($1, $2, $3, $4)
		RETURNING id, translation, chunks, pinyin, hsk_level, created_at, updated_at
	`, req.Translation, req.Chunks, req.Pinyin, req.HSKLevel).Scan(
		&e.ID, &e.Translation, &e.Chunks, &e.Pinyin, &e.HSKLevel, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req AdminExerciseRequest) (*Exercise, error) {

	e := &Exercise{}

	err := r.db.QueryRow(ctx, `
		UPDATE sentence_exercises
		SET translation = $1, chunks = $2, pinyin = $3, hsk_level = $4, updated_at = now()
		WHERE id = $5
		RETURNING id, translation, chunks, pinyin, hsk_level, created_at, updated_at
	`, req.Translation, req.Chunks, req.Pinyin, req.HSKLevel, id).Scan(
		&e.ID, &e.Translation, &e.Chunks, &e.Pinyin, &e.HSKLevel, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sentence_exercises WHERE id = $1`, id)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Exercise, error) {

	var e Exercise

	err := r.db.QueryRow(ctx, `
		SELECT id, translation, chunks, pinyin, hsk_level, created_at, updated_at
		FROM sentence_exercises
		WHERE id = $1
	`, id).Scan(
		&e.ID,
		&e.Translation,
		&e.Chunks,
		&e.Pinyin,
		&e.HSKLevel,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int64) ([]Exercise, error) {

	if len(ids) == 0 {
		return []Exercise{}, nil
	}

	args := make([]interface{}, len(ids))
	placeholders := make([]string, len(ids))

	for i, id := range ids {
		args[i] = id
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}

	query := `
		SELECT id, translation, chunks, pinyin, hsk_level, created_at, updated_at
		FROM sentence_exercises
		WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Exercise, 0, len(ids))

	for rows.Next() {

		var e Exercise

		if err := rows.Scan(
			&e.ID,
			&e.Translation,
			&e.Chunks,
			&e.Pinyin,
			&e.HSKLevel,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, rows.Err()
}
