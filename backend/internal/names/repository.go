package names

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

func (r *Repository) GetAll(ctx context.Context) ([]Name, error) {

	rows, err := r.db.Query(ctx, `
		SELECT id, hanzi, pinyin, translation, created_at, updated_at
		FROM names
		ORDER BY hanzi
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]Name, 0)

	for rows.Next() {

		var n Name

		if err := rows.Scan(
			&n.ID,
			&n.Hanzi,
			&n.Pinyin,
			&n.Translation,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, n)
	}

	return result, rows.Err()
}
