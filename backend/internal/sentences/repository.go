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
