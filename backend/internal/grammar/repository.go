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
			COALESCE(example2_hanzi, ''),
			COALESCE(example2_pinyin, ''),
			COALESCE(example2_translation, ''),
			COALESCE(example3_hanzi, ''),
			COALESCE(example3_pinyin, ''),
			COALESCE(example3_translation, ''),
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
			&n.Example2Hanzi,
			&n.Example2Pinyin,
			&n.Example2Translation,
			&n.Example3Hanzi,
			&n.Example3Pinyin,
			&n.Example3Translation,
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

func (r *Repository) AdminList(ctx context.Context, request AdminListRequest) ([]Note, int, error) {

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
				title ILIKE $` + strconv.Itoa(arg) + `
				OR explanation ILIKE $` + strconv.Itoa(arg) + `
			)
		`

		args = append(args, "%"+request.Search+"%")
		arg++
	}

	if request.HSK > 0 {
		where += ` AND hsk_level = $` + strconv.Itoa(arg)

		args = append(args, request.HSK)
		arg++
	}

	var total int

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM grammar_notes `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id,
			title,
			explanation,
			COALESCE(example_hanzi, ''),
			COALESCE(example_pinyin, ''),
			COALESCE(example_translation, ''),
			COALESCE(example2_hanzi, ''),
			COALESCE(example2_pinyin, ''),
			COALESCE(example2_translation, ''),
			COALESCE(example3_hanzi, ''),
			COALESCE(example3_pinyin, ''),
			COALESCE(example3_translation, ''),
			hsk_level,
			created_at,
			updated_at
		FROM grammar_notes
	` + where + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	result := make([]Note, 0)

	for rows.Next() {

		var n Note

		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Explanation,
			&n.ExampleHanzi,
			&n.ExamplePinyin,
			&n.ExampleTranslation,
			&n.Example2Hanzi,
			&n.Example2Pinyin,
			&n.Example2Translation,
			&n.Example3Hanzi,
			&n.Example3Pinyin,
			&n.Example3Translation,
			&n.HSKLevel,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		result = append(result, n)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *Repository) Create(ctx context.Context, req CreateNoteRequest) (*Note, error) {

	n := &Note{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO grammar_notes (
			title, explanation,
			example_hanzi, example_pinyin, example_translation,
			example2_hanzi, example2_pinyin, example2_translation,
			example3_hanzi, example3_pinyin, example3_translation,
			hsk_level
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING
			id, title, explanation,
			COALESCE(example_hanzi, ''), COALESCE(example_pinyin, ''), COALESCE(example_translation, ''),
			COALESCE(example2_hanzi, ''), COALESCE(example2_pinyin, ''), COALESCE(example2_translation, ''),
			COALESCE(example3_hanzi, ''), COALESCE(example3_pinyin, ''), COALESCE(example3_translation, ''),
			hsk_level, created_at, updated_at
	`,
		req.Title, req.Explanation,
		req.ExampleHanzi, req.ExamplePinyin, req.ExampleTranslation,
		req.Example2Hanzi, req.Example2Pinyin, req.Example2Translation,
		req.Example3Hanzi, req.Example3Pinyin, req.Example3Translation,
		req.HSKLevel,
	).Scan(
		&n.ID, &n.Title, &n.Explanation,
		&n.ExampleHanzi, &n.ExamplePinyin, &n.ExampleTranslation,
		&n.Example2Hanzi, &n.Example2Pinyin, &n.Example2Translation,
		&n.Example3Hanzi, &n.Example3Pinyin, &n.Example3Translation,
		&n.HSKLevel, &n.CreatedAt, &n.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req UpdateNoteRequest) (*Note, error) {

	n := &Note{}

	err := r.db.QueryRow(ctx, `
		UPDATE grammar_notes
		SET
			title = $1, explanation = $2,
			example_hanzi = $3, example_pinyin = $4, example_translation = $5,
			example2_hanzi = $6, example2_pinyin = $7, example2_translation = $8,
			example3_hanzi = $9, example3_pinyin = $10, example3_translation = $11,
			hsk_level = $12, updated_at = now()
		WHERE id = $13
		RETURNING
			id, title, explanation,
			COALESCE(example_hanzi, ''), COALESCE(example_pinyin, ''), COALESCE(example_translation, ''),
			COALESCE(example2_hanzi, ''), COALESCE(example2_pinyin, ''), COALESCE(example2_translation, ''),
			COALESCE(example3_hanzi, ''), COALESCE(example3_pinyin, ''), COALESCE(example3_translation, ''),
			hsk_level, created_at, updated_at
	`,
		req.Title, req.Explanation,
		req.ExampleHanzi, req.ExamplePinyin, req.ExampleTranslation,
		req.Example2Hanzi, req.Example2Pinyin, req.Example2Translation,
		req.Example3Hanzi, req.Example3Pinyin, req.Example3Translation,
		req.HSKLevel, id,
	).Scan(
		&n.ID, &n.Title, &n.Explanation,
		&n.ExampleHanzi, &n.ExamplePinyin, &n.ExampleTranslation,
		&n.Example2Hanzi, &n.Example2Pinyin, &n.Example2Translation,
		&n.Example3Hanzi, &n.Example3Pinyin, &n.Example3Translation,
		&n.HSKLevel, &n.CreatedAt, &n.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM grammar_notes WHERE id = $1`, id)
	return err
}
