package learning

import (
	"context"
	"errors"

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

func (r *Repository) Get(ctx context.Context, userID, wordID int64) (*WordLearningProgress, error) {

	var progress WordLearningProgress

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			word_id,
			stage,
			next_eligible_at,
			learned_at,
			created_at,
			updated_at
		FROM word_learning_progress
		WHERE user_id = $1
		  AND word_id = $2
	`, userID, wordID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.WordID,
		&progress.Stage,
		&progress.NextEligibleAt,
		&progress.LearnedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &progress, nil
}

func (r *Repository) Upsert(ctx context.Context, progress *WordLearningProgress) error {

	return r.db.QueryRow(ctx, `
		INSERT INTO word_learning_progress (
			user_id,
			word_id,
			stage,
			next_eligible_at,
			learned_at
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, word_id)
		DO UPDATE SET
			stage = EXCLUDED.stage,
			next_eligible_at = EXCLUDED.next_eligible_at,
			learned_at = EXCLUDED.learned_at,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`,
		progress.UserID,
		progress.WordID,
		progress.Stage,
		progress.NextEligibleAt,
		progress.LearnedAt,
	).Scan(&progress.ID, &progress.CreatedAt, &progress.UpdatedAt)
}

func (r *Repository) ListForUser(ctx context.Context, userID int64) ([]WordLearningProgress, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			user_id,
			word_id,
			stage,
			next_eligible_at,
			learned_at,
			created_at,
			updated_at
		FROM word_learning_progress
		WHERE user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]WordLearningProgress, 0)

	for rows.Next() {

		var progress WordLearningProgress

		if err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.WordID,
			&progress.Stage,
			&progress.NextEligibleAt,
			&progress.LearnedAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, progress)
	}

	return items, rows.Err()
}

func (r *Repository) ListInProgress(ctx context.Context, userID int64) ([]WordProgressDetail, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			w.id,
			w.hanzi,
			w.pinyin,
			w.translation,
			w.hsk_level,
			p.stage,
			p.next_eligible_at
		FROM word_learning_progress p
		INNER JOIN words w
			ON w.id = p.word_id
		WHERE p.user_id = $1
		  AND p.learned_at IS NULL
		  AND p.stage > 0
		ORDER BY p.next_eligible_at ASC NULLS LAST
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]WordProgressDetail, 0)

	for rows.Next() {

		var detail WordProgressDetail

		if err := rows.Scan(
			&detail.WordID,
			&detail.Hanzi,
			&detail.Pinyin,
			&detail.Translation,
			&detail.HSKLevel,
			&detail.Stage,
			&detail.NextEligibleAt,
		); err != nil {
			return nil, err
		}

		result = append(result, detail)
	}

	return result, rows.Err()
}

func (r *Repository) ListLearnedWords(ctx context.Context, userID int64) ([]words.Word, error) {

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
		FROM word_learning_progress p
		INNER JOIN words w
			ON w.id = p.word_id
		WHERE p.user_id = $1
		  AND p.learned_at IS NOT NULL
		ORDER BY p.learned_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]words.Word, 0)

	for rows.Next() {

		var word words.Word

		if err := rows.Scan(
			&word.ID,
			&word.Hanzi,
			&word.Pinyin,
			&word.Translation,
			&word.PartOfSpeech,
			&word.HSKLevel,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, word)
	}

	return result, rows.Err()
}
