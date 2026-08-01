package texts

import (
	"context"
	"strconv"

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
			t.id, t.title, t.hanzi, t.pinyin, t.translation, t.hsk_level, t.image_url,
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
			&t.ImageURL,
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
			t.id, t.title, t.hanzi, t.pinyin, t.translation, t.hsk_level, t.image_url,
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
		&t.ImageURL,
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

// AdminList is independent of List: no per-user progress join, has a total
// count for pagination, and orders newest-first so new content is easy to
// find instead of being buried by HSK level.
func (r *Repository) AdminList(ctx context.Context, request AdminListRequest) ([]Text, int, error) {

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
				OR hanzi ILIKE $` + strconv.Itoa(arg) + `
				OR translation ILIKE $` + strconv.Itoa(arg) + `
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

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM texts `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, title, hanzi, pinyin, translation, hsk_level, image_url, '' AS status, created_at, updated_at
		FROM texts
	` + where + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
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
			&t.ImageURL,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *Repository) Create(ctx context.Context, req CreateTextRequest) (*Text, error) {

	t := &Text{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO texts (title, hanzi, pinyin, translation, hsk_level, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, hanzi, pinyin, translation, hsk_level, image_url, created_at, updated_at
	`, req.Title, req.Hanzi, req.Pinyin, req.Translation, req.HSKLevel, req.ImageURL).Scan(
		&t.ID, &t.Title, &t.Hanzi, &t.Pinyin, &t.Translation, &t.HSKLevel, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req UpdateTextRequest) (*Text, error) {

	t := &Text{}

	err := r.db.QueryRow(ctx, `
		UPDATE texts
		SET title = $1, hanzi = $2, pinyin = $3, translation = $4, hsk_level = $5, image_url = $6, updated_at = now()
		WHERE id = $7
		RETURNING id, title, hanzi, pinyin, translation, hsk_level, image_url, created_at, updated_at
	`, req.Title, req.Hanzi, req.Pinyin, req.Translation, req.HSKLevel, req.ImageURL, id).Scan(
		&t.ID, &t.Title, &t.Hanzi, &t.Pinyin, &t.Translation, &t.HSKLevel, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM texts WHERE id = $1`, id)
	return err
}
