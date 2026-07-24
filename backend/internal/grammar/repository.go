package grammar

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

func (r *Repository) GetByIDs(ctx context.Context, ids []int64) ([]Note, error) {

	if len(ids) == 0 {
		return []Note{}, nil
	}

	args := make([]interface{}, len(ids))
	placeholders := make([]string, len(ids))

	for i, id := range ids {
		args[i] = id
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}

	query := `
		SELECT
			id,
			title,
			explanation,
			COALESCE(example_hanzi, ''),
			COALESCE(example_pinyin, ''),
			COALESCE(example_translation, ''),
			hsk_level,
			created_at,
			updated_at
		FROM grammar_notes
		WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Note, 0, len(ids))

	for rows.Next() {

		var n Note

		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Explanation,
			&n.ExampleHanzi,
			&n.ExamplePinyin,
			&n.ExampleTranslation,
			&n.HSKLevel,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, n)
	}

	return result, rows.Err()
}
