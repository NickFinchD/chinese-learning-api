package words

import (
	"context"
	"strconv"
	"strings"

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

	var words []Word

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

	return word, nil
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

	var words []Word

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
