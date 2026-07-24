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

func (r *Repository) List(ctx context.Context, hskLevel int16, userID int64) ([]Text, error) {

	query := `
		SELECT
			t.id, t.title, t.hanzi, t.pinyin, t.translation, t.hsk_level,
			COALESCE(utp.status, 'not_started') AS status,
			t.created_at, t.updated_at
		FROM texts t
		LEFT JOIN user_text_progress utp ON utp.text_id = t.id AND utp.user_id = $1
	`

	args := []interface{}{userID}

	if hskLevel > 0 {
		query += ` WHERE t.hsk_level = $2`
		args = append(args, hskLevel)
	}

	query += ` ORDER BY t.hsk_level, t.id`

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
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id int64, userID int64) (*Text, error) {

	var t Text

	err := r.db.QueryRow(ctx, `
		SELECT
			t.id, t.title, t.hanzi, t.pinyin, t.translation, t.hsk_level,
			COALESCE(utp.status, 'not_started') AS status,
			t.created_at, t.updated_at
		FROM texts t
		LEFT JOIN user_text_progress utp ON utp.text_id = t.id AND utp.user_id = $2
		WHERE t.id = $1
	`, id, userID).Scan(
		&t.ID,
		&t.Title,
		&t.Hanzi,
		&t.Pinyin,
		&t.Translation,
		&t.HSKLevel,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// MarkStarted records that the user has opened this text, unless it's
// already tracked (so re-opening a completed text doesn't downgrade it).
func (r *Repository) MarkStarted(ctx context.Context, userID, textID int64) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_text_progress (user_id, text_id, status)
		VALUES ($1, $2, 'in_progress')
		ON CONFLICT (user_id, text_id) DO NOTHING
	`, userID, textID)

	return err
}

func (r *Repository) MarkRead(ctx context.Context, userID, textID int64) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_text_progress (user_id, text_id, status, read_at)
		VALUES ($1, $2, 'completed', NOW())
		ON CONFLICT (user_id, text_id)
		DO UPDATE SET status = 'completed', read_at = NOW(), updated_at = NOW()
	`, userID, textID)

	return err
}

func (r *Repository) MarkUnread(ctx context.Context, userID, textID int64) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_text_progress (user_id, text_id, status, read_at)
		VALUES ($1, $2, 'in_progress', NULL)
		ON CONFLICT (user_id, text_id)
		DO UPDATE SET status = 'in_progress', read_at = NULL, updated_at = NOW()
	`, userID, textID)

	return err
}
