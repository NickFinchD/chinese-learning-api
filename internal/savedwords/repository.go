package savedwords

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
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

func (r *Repository) Save(ctx context.Context, userID, wordID int64) error {

	query := `
		INSERT INTO saved_words (
			user_id,
			word_id
		)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, userID, wordID)

	return err
}
func (r *Repository) List(ctx context.Context, userID int64) ([]words.Word, error) {

	query := `
		SELECT
			w.id,
			w.hanzi,
			w.pinyin,
			w.translation,
			w.part_of_speech,
			w.hsk_level,
			w.created_at,
			w.updated_at
		FROM saved_words sw
		INNER JOIN words w
			ON w.id = sw.word_id
		WHERE sw.user_id = $1
		ORDER BY sw.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []words.Word

	for rows.Next() {

		var word words.Word

		err := rows.Scan(
			&word.ID,
			&word.Hanzi,
			&word.Pinyin,
			&word.Translation,
			&word.PartOfSpeech,
			&word.HSKLevel,
			&word.CreatedAt,
			&word.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, word)
	}

	return result, rows.Err()
}
func (r *Repository) Delete(ctx context.Context, userID, wordID int64) error {

	query := `
		DELETE FROM saved_words
		WHERE user_id = $1
		  AND word_id = $2
	`

	_, err := r.db.Exec(ctx, query, userID, wordID)

	return err
}
