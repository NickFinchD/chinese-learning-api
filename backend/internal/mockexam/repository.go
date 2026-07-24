package mockexam

import (
	"context"

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

func (r *Repository) CreateAttempt(
	ctx context.Context,
	userID int64,
	hskLevel int16,
	totalQuestions, correctCount, scorePercent int16,
	passed bool,
	durationSeconds int32,
) (*Attempt, error) {

	a := Attempt{
		UserID:          userID,
		HSKLevel:        hskLevel,
		TotalQuestions:  totalQuestions,
		CorrectCount:    correctCount,
		ScorePercent:    scorePercent,
		Passed:          passed,
		DurationSeconds: durationSeconds,
	}

	err := r.db.QueryRow(ctx, `
		INSERT INTO mock_exam_attempts (
			user_id, hsk_level, total_questions, correct_count,
			score_percent, passed, duration_seconds
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`,
		userID, hskLevel, totalQuestions, correctCount,
		scorePercent, passed, durationSeconds,
	).Scan(&a.ID, &a.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *Repository) HasPassed(ctx context.Context, userID int64, hskLevel int16) (bool, error) {

	var exists bool

	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM mock_exam_attempts
			WHERE user_id = $1 AND hsk_level = $2 AND passed = TRUE
		)
	`, userID, hskLevel).Scan(&exists)

	return exists, err
}

func (r *Repository) ListAttempts(ctx context.Context, userID int64) ([]Attempt, error) {

	rows, err := r.db.Query(ctx, `
		SELECT id, hsk_level, total_questions, correct_count,
		       score_percent, passed, duration_seconds, created_at
		FROM mock_exam_attempts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 20
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]Attempt, 0)

	for rows.Next() {

		var a Attempt

		if err := rows.Scan(
			&a.ID, &a.HSKLevel, &a.TotalQuestions, &a.CorrectCount,
			&a.ScorePercent, &a.Passed, &a.DurationSeconds, &a.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, a)
	}

	return result, rows.Err()
}
