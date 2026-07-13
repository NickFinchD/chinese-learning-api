package review

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

func (r *Repository) GetWordsForReview(
	ctx context.Context,
	userID int64,
) ([]ReviewWordResponse, error) {

	query := `
		SELECT
			w.id,
			w.hanzi,
			w.pinyin,
			w.translation
		FROM user_word_progress uwp
		INNER JOIN words w
			ON w.id = uwp.word_id
		WHERE
			uwp.user_id = $1
			AND uwp.next_review_at <= NOW()
		ORDER BY uwp.next_review_at
	`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var words []ReviewWordResponse

	for rows.Next() {

		var word ReviewWordResponse

		err := rows.Scan(
			&word.WordID,
			&word.Hanzi,
			&word.Pinyin,
			&word.Translation,
		)

		if err != nil {
			return nil, err
		}

		words = append(words, word)
	}

	return words, rows.Err()
}
