package collections

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) Create(ctx context.Context, userID int64, name string) (*Collection, error) {

	c := Collection{UserID: userID, Name: name}

	err := r.db.QueryRow(ctx, `
		INSERT INTO word_collections (user_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`, userID, name).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) List(ctx context.Context, userID int64) ([]Collection, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			wc.id,
			wc.name,
			wc.created_at,
			wc.updated_at,
			COUNT(wci.id) AS word_count
		FROM word_collections wc
		LEFT JOIN word_collection_items wci ON wci.collection_id = wc.id
		WHERE wc.user_id = $1
		GROUP BY wc.id
		ORDER BY wc.created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]Collection, 0)

	for rows.Next() {

		var c Collection

		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt, &c.WordCount); err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, userID, collectionID int64) (*Collection, error) {

	var c Collection

	err := r.db.QueryRow(ctx, `
		SELECT
			wc.id,
			wc.name,
			wc.created_at,
			wc.updated_at,
			COUNT(wci.id) AS word_count
		FROM word_collections wc
		LEFT JOIN word_collection_items wci ON wci.collection_id = wc.id
		WHERE wc.user_id = $1 AND wc.id = $2
		GROUP BY wc.id
	`, userID, collectionID).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt, &c.WordCount)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *Repository) ListWords(ctx context.Context, collectionID int64) ([]words.Word, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			w.id,
			w.hanzi,
			w.pinyin,
			w.translation,
			w.part_of_speech,
			w.hsk_level,
			w.created_at,
			w.updated_at
		FROM word_collection_items wci
		INNER JOIN words w ON w.id = wci.word_id
		WHERE wci.collection_id = $1
		ORDER BY wci.created_at DESC
	`, collectionID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]words.Word, 0)

	for rows.Next() {

		var w words.Word

		if err := rows.Scan(
			&w.ID,
			&w.Hanzi,
			&w.Pinyin,
			&w.Translation,
			&w.PartOfSpeech,
			&w.HSKLevel,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, w)
	}

	return result, rows.Err()
}

func (r *Repository) Rename(ctx context.Context, userID, collectionID int64, name string) error {

	tag, err := r.db.Exec(ctx, `
		UPDATE word_collections
		SET name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
	`, name, collectionID, userID)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, userID, collectionID int64) error {

	_, err := r.db.Exec(ctx, `
		DELETE FROM word_collections
		WHERE id = $1 AND user_id = $2
	`, collectionID, userID)

	return err
}

// AddWord is a no-op (ON CONFLICT DO NOTHING) if the word is already in the
// collection, and silently affects zero rows if the collection isn't owned
// by userID — callers that need to distinguish "not found" should check
// ownership via GetByID first.
func (r *Repository) AddWord(ctx context.Context, userID, collectionID, wordID int64) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO word_collection_items (collection_id, word_id)
		SELECT $1, $2
		WHERE EXISTS (
			SELECT 1 FROM word_collections WHERE id = $1 AND user_id = $3
		)
		ON CONFLICT (collection_id, word_id) DO NOTHING
	`, collectionID, wordID, userID)

	return err
}

func (r *Repository) RemoveWord(ctx context.Context, userID, collectionID, wordID int64) error {

	_, err := r.db.Exec(ctx, `
		DELETE FROM word_collection_items
		WHERE collection_id = $1
		  AND word_id = $2
		  AND EXISTS (
			SELECT 1 FROM word_collections WHERE id = $1 AND user_id = $3
		  )
	`, collectionID, wordID, userID)

	return err
}
