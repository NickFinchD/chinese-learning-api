package quizzes

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

func (r *Repository) GetByIDs(ctx context.Context, ids []int64) ([]Quiz, error) {

	if len(ids) == 0 {
		return []Quiz{}, nil
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
			question,
			hsk_level,
			created_at,
			updated_at
		FROM quizzes
		WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	quizzes := make([]Quiz, 0)

	for rows.Next() {

		var quiz Quiz

		err := rows.Scan(
			&quiz.ID,
			&quiz.Question,
			&quiz.HSKLevel,
			&quiz.CreatedAt,
			&quiz.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		quizzes = append(quizzes, quiz)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.loadOptions(ctx, quizzes); err != nil {
		return nil, err
	}

	return quizzes, nil
}
func (r *Repository) loadOptions(ctx context.Context, quizzes []Quiz) error {

	if len(quizzes) == 0 {
		return nil
	}

	ids := make([]interface{}, len(quizzes))
	placeholders := make([]string, len(quizzes))

	for i, quiz := range quizzes {
		ids[i] = quiz.ID
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}

	query := `
		SELECT
			id,
			quiz_id,
			option_text,
			is_correct,
			sort_order
		FROM quiz_options
		WHERE quiz_id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY quiz_id, sort_order
	`

	rows, err := r.db.Query(ctx, query, ids...)
	if err != nil {
		return err
	}

	defer rows.Close()

	options := make(map[int64][]Option)

	for rows.Next() {

		var option Option

		err := rows.Scan(
			&option.ID,
			&option.QuizID,
			&option.Text,
			&option.IsCorrect,
			&option.SortOrder,
		)

		if err != nil {
			return err
		}

		options[option.QuizID] = append(options[option.QuizID], option)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for i := range quizzes {
		quizzes[i].Options = options[quizzes[i].ID]
	}

	return nil
}
func (r *Repository) GetAll(ctx context.Context) ([]Quiz, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, question, hsk_level, created_at, updated_at
		FROM quizzes
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quizzes := make([]Quiz, 0)

	for rows.Next() {
		var quiz Quiz

		if err := rows.Scan(
			&quiz.ID,
			&quiz.Question,
			&quiz.HSKLevel,
			&quiz.CreatedAt,
			&quiz.UpdatedAt,
		); err != nil {
			return nil, err
		}

		quizzes = append(quizzes, quiz)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.loadOptions(ctx, quizzes); err != nil {
		return nil, err
	}

	return quizzes, nil
}
func (r *Repository) GetByHSKLevel(ctx context.Context, hsk int16) ([]Quiz, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, question, hsk_level, created_at, updated_at
		FROM quizzes
		WHERE hsk_level = $1
		ORDER BY id
	`, hsk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quizzes := make([]Quiz, 0)

	for rows.Next() {
		var quiz Quiz

		if err := rows.Scan(
			&quiz.ID,
			&quiz.Question,
			&quiz.HSKLevel,
			&quiz.CreatedAt,
			&quiz.UpdatedAt,
		); err != nil {
			return nil, err
		}

		quizzes = append(quizzes, quiz)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.loadOptions(ctx, quizzes); err != nil {
		return nil, err
	}

	return quizzes, nil
}
func (r *Repository) GetByID(ctx context.Context, id int64) (*Quiz, error) {
	var quiz Quiz

	err := r.db.QueryRow(ctx, `
		SELECT id, question, hsk_level, created_at, updated_at
		FROM quizzes
		WHERE id = $1
	`, id).Scan(
		&quiz.ID,
		&quiz.Question,
		&quiz.HSKLevel,
		&quiz.CreatedAt,
		&quiz.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	quizzes := []Quiz{quiz}

	if err := r.loadOptions(ctx, quizzes); err != nil {
		return nil, err
	}

	return &quizzes[0], nil
}
func (r *Repository) Create(ctx context.Context, quiz Quiz) (*Quiz, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO quizzes (question)
		VALUES ($1)
		RETURNING id, hsk_level, created_at, updated_at
	`, quiz.Question).Scan(
		&quiz.ID,
		&quiz.HSKLevel,
		&quiz.CreatedAt,
		&quiz.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	for i, option := range quiz.Options {
		err = tx.QueryRow(ctx, `
			INSERT INTO quiz_options (
				quiz_id,
				option_text,
				is_correct,
				sort_order
			)
			VALUES ($1,$2,$3,$4)
			RETURNING id
		`,
			quiz.ID,
			option.Text,
			option.IsCorrect,
			i+1,
		).Scan(&option.ID)

		if err != nil {
			return nil, err
		}

		option.QuizID = quiz.ID
		option.SortOrder = int16(i + 1)

		quiz.Options[i] = option
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &quiz, nil
}
func (r *Repository) CheckAnswer(ctx context.Context, quizID, optionID int64) (bool, error) {
	var correct bool

	err := r.db.QueryRow(ctx, `
		SELECT is_correct
		FROM quiz_options
		WHERE id = $1
		  AND quiz_id = $2
	`,
		optionID,
		quizID,
	).Scan(&correct)

	if err != nil {
		return false, err
	}

	return correct, nil
}
