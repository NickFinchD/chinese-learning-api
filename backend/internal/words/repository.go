package words

import (
	"context"
	"strconv"
	"strings"

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

func (r *Repository) List(ctx context.Context, request ListRequest) ([]Word, error) {

	query := `
		SELECT
			id,
			hanzi,
			pinyin,
			translation,
			part_of_speech,
			hsk_level,
			created_at,
			updated_at
		FROM words
		WHERE 1=1
	`

	args := []interface{}{}
	arg := 1

	if request.Search != "" {
		query += `
			AND (
				hanzi ILIKE $` + strconv.Itoa(arg) + `
				OR pinyin ILIKE $` + strconv.Itoa(arg) + `
				OR translation ILIKE $` + strconv.Itoa(arg) + `
			)
		`

		args = append(args, "%"+request.Search+"%")
		arg++
	}

	if request.HSK > 0 {
		query += `
			AND hsk_level = $` + strconv.Itoa(arg)

		args = append(args, request.HSK)
		arg++
	}

	query += `
		ORDER BY hsk_level, id
	`

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	words := make([]Word, 0)

	for rows.Next() {

		var word Word

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

		words = append(words, word)
	}

	return words, rows.Err()
}
func (r *Repository) GetByID(ctx context.Context, id int64) (*Word, error) {

	query := `
		SELECT
			id,
			hanzi,
			pinyin,
			translation,
			part_of_speech,
			hsk_level,
			created_at,
			updated_at
		FROM words
		WHERE id = $1
	`

	word := &Word{}

	err := r.db.QueryRow(ctx, query, id).Scan(
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

	examples, err := r.GetExamplesByWordID(ctx, id)
	if err != nil {
		return nil, err
	}

	word.Examples = examples

	return word, nil
}

func (r *Repository) GetExamplesByWordID(ctx context.Context, wordID int64) ([]Example, error) {

	rows, err := r.db.Query(ctx, `
		SELECT id, hanzi, pinyin, translation
		FROM word_examples
		WHERE word_id = $1
		ORDER BY sort_order
	`, wordID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	examples := make([]Example, 0)

	for rows.Next() {

		var e Example

		if err := rows.Scan(&e.ID, &e.Hanzi, &e.Pinyin, &e.Translation); err != nil {
			return nil, err
		}

		examples = append(examples, e)
	}

	return examples, rows.Err()
}

// AdminList is a separate query from List because the admin screen has
// different needs: a total count for pagination, and no join with any
// per-user state (List has none either today, but admin listing is kept
// independent so it doesn't pick up user-facing filtering assumptions).
func (r *Repository) AdminList(ctx context.Context, request AdminListRequest) ([]Word, int, error) {

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
				hanzi ILIKE $` + strconv.Itoa(arg) + `
				OR pinyin ILIKE $` + strconv.Itoa(arg) + `
				OR translation ILIKE $` + strconv.Itoa(arg) + `
			)
		`

		args = append(args, "%"+request.Search+"%")
		arg++
	}

	if request.HSK > 0 {
		where += `
			AND hsk_level = $` + strconv.Itoa(arg)

		args = append(args, request.HSK)
		arg++
	}

	var total int

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM words `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id,
			hanzi,
			pinyin,
			translation,
			part_of_speech,
			hsk_level,
			created_at,
			updated_at
		FROM words
	` + where + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	words := make([]Word, 0)

	for rows.Next() {

		var word Word

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
			return nil, 0, err
		}

		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return words, total, nil
}

func (r *Repository) Create(ctx context.Context, req CreateWordRequest) (*Word, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	word := &Word{}

	err = tx.QueryRow(ctx, `
		INSERT INTO words (hanzi, pinyin, translation, part_of_speech, hsk_level)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, hanzi, pinyin, translation, part_of_speech, hsk_level, created_at, updated_at
	`, req.Hanzi, req.Pinyin, req.Translation, req.PartOfSpeech, req.HSKLevel).Scan(
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

	if err := insertExamples(ctx, tx, word.ID, req.Examples); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, word.ID)
}

// Update replaces a word's examples wholesale (delete then reinsert) rather
// than diffing — the simplest correct approach at this scale, matching how
// quiz options are handled.
func (r *Repository) Update(ctx context.Context, id int64, req UpdateWordRequest) (*Word, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE words
		SET hanzi = $1, pinyin = $2, translation = $3, part_of_speech = $4, hsk_level = $5, updated_at = now()
		WHERE id = $6
	`, req.Hanzi, req.Pinyin, req.Translation, req.PartOfSpeech, req.HSKLevel, id)

	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM word_examples WHERE word_id = $1`, id); err != nil {
		return nil, err
	}

	if err := insertExamples(ctx, tx, id, req.Examples); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func insertExamples(ctx context.Context, tx pgx.Tx, wordID int64, examples []ExampleRequest) error {

	for i, example := range examples {

		_, err := tx.Exec(ctx, `
			INSERT INTO word_examples (word_id, hanzi, pinyin, translation, sort_order)
			VALUES ($1, $2, $3, $4, $5)
		`, wordID, example.Hanzi, example.Pinyin, example.Translation, i+1)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM words WHERE id = $1`, id)
	return err
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int64) ([]Word, error) {

	if len(ids) == 0 {
		return []Word{}, nil
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
			hanzi,
			pinyin,
			translation,
			part_of_speech,
			hsk_level,
			created_at,
			updated_at
		FROM words
		WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	words := make([]Word, 0)

	for rows.Next() {

		var word Word

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

		words = append(words, word)
	}

	return words, rows.Err()
}
