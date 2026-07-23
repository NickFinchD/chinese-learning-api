package texts

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

func (r *Repository) List(ctx context.Context, hskLevel int16) ([]Text, error) {

	query := `
		SELECT id, title, hanzi, pinyin, translation, hsk_level, created_at, updated_at
		FROM texts
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

	result := make([]Text, 0)

	for rows.Next() {

		var t Text

		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Hanzi,
			&t.Pinyin,
			&t.Translation,
			&t.HSKLevel,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Text, error) {

	var t Text

	err := r.db.QueryRow(ctx, `
		SELECT id, title, hanzi, pinyin, translation, hsk_level, created_at, updated_at
		FROM texts
		WHERE id = $1
	`, id).Scan(
		&t.ID,
		&t.Title,
		&t.Hanzi,
		&t.Pinyin,
		&t.Translation,
		&t.HSKLevel,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
