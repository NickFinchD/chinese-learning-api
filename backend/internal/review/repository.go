package review

import (
	"context"
	"errors"
	"time"

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

	words := make([]ReviewWordResponse, 0)

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
func (r *Repository) AddWord(
	ctx context.Context,
	userID int64,
	wordID int64,
) error {

	query := `
		INSERT INTO user_word_progress (
			user_id,
			word_id
		)
		VALUES (
			$1,
			$2
		)
		ON CONFLICT (user_id, word_id)
		DO NOTHING
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		wordID,
	)

	return err
}
func (r *Repository) UpdateProgress(
	ctx context.Context,
	userID int64,
	wordID int64,
	reviewCount int,
	correctCount int,
	wrongCount int,
	nextReviewAt time.Time,
	lastReviewAt time.Time,
) error {

	query := `
		UPDATE user_word_progress
		SET
			review_count = $3,
			correct_count = $4,
			wrong_count = $5,
			next_review_at = $6,
			last_review_at = $7,
			updated_at = NOW()
		WHERE
			user_id = $1
			AND word_id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		wordID,
		reviewCount,
		correctCount,
		wrongCount,
		nextReviewAt,
		lastReviewAt,
	)

	return err
}
func (r *Repository) GetProgress(
	ctx context.Context,
	userID int64,
	wordID int64,
) (*UserWordProgress, error) {

	query := `
		SELECT
			id,
			user_id,
			word_id,
			review_count,
			correct_count,
			wrong_count,
			ease_factor,
			last_review_at,
			next_review_at,
			created_at,
			updated_at
		FROM user_word_progress
		WHERE
			user_id = $1
			AND word_id = $2
	`

	progress := &UserWordProgress{}

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		wordID,
	).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.WordID,
		&progress.ReviewCount,
		&progress.CorrectCount,
		&progress.WrongCount,
		&progress.EaseFactor,
		&progress.LastReviewAt,
		&progress.NextReviewAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWordNotFound
		}

		return nil, err
	}

	return progress, nil
}
func (r *Repository) GetStatistics(
	ctx context.Context,
	userID int64,
) (*StatisticsResponse, error) {

	stats := &StatisticsResponse{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM user_word_progress
		WHERE user_id = $1
		`,
		userID,
	).Scan(&stats.TotalWords)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM user_word_progress
		WHERE
			user_id = $1
			AND next_review_at <= NOW()
		`,
		userID,
	).Scan(&stats.ReadyForReview)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM user_word_progress
		WHERE
			user_id = $1
			AND review_count > 0
		`,
		userID,
	).Scan(&stats.ReviewedWords)

	if err != nil {
		return nil, err
	}

	return stats, nil
}
