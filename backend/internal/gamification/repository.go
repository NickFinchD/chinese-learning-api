package gamification

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

func (r *Repository) GetStats(ctx context.Context, userID int64) (*UserStats, error) {

	var stats UserStats

	err := r.db.QueryRow(ctx, `
		SELECT user_id, xp, total_seconds_active, last_heartbeat_at
		FROM user_stats
		WHERE user_id = $1
	`, userID).Scan(
		&stats.UserID,
		&stats.XP,
		&stats.TotalSecondsActive,
		&stats.LastHeartbeatAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &UserStats{UserID: userID}, nil
		}

		return nil, err
	}

	return &stats, nil
}

// Heartbeat records that the user is currently active. The amount of time
// added to their total is the real wall-clock gap since their previous
// heartbeat, capped at maxGapSeconds so a stale tab or a sleeping machine
// can't inflate the total after a long absence.
func (r *Repository) Heartbeat(ctx context.Context, userID int64, maxGapSeconds int64) error {

	_, err := r.db.Exec(ctx, `
		WITH prev AS (
			SELECT last_heartbeat_at FROM user_stats WHERE user_id = $1
		)
		INSERT INTO user_stats (user_id, total_seconds_active, last_heartbeat_at, updated_at)
		VALUES ($1, 0, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			total_seconds_active = user_stats.total_seconds_active + LEAST(
				$2::bigint,
				GREATEST(0, EXTRACT(EPOCH FROM (NOW() - COALESCE((SELECT last_heartbeat_at FROM prev), NOW())))::bigint)
			),
			last_heartbeat_at = NOW(),
			updated_at = NOW()
	`, userID, maxGapSeconds)

	return err
}

func (r *Repository) AddXP(ctx context.Context, userID int64, amount int) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_stats (user_id, xp, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			xp = user_stats.xp + $2,
			updated_at = NOW()
	`, userID, amount)

	return err
}

func (r *Repository) ListAchievements(ctx context.Context) ([]Achievement, error) {

	rows, err := r.db.Query(ctx, `
		SELECT id, code, title, description, tier, metric, threshold, xp_reward
		FROM achievements
		ORDER BY metric, tier
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Achievement, 0)

	for rows.Next() {

		var a Achievement

		if err := rows.Scan(
			&a.ID, &a.Code, &a.Title, &a.Description, &a.Tier, &a.Metric, &a.Threshold, &a.XPReward,
		); err != nil {
			return nil, err
		}

		result = append(result, a)
	}

	return result, rows.Err()
}

func (r *Repository) ListUnlockedAchievementIDs(ctx context.Context, userID int64) (map[int64]time.Time, error) {

	rows, err := r.db.Query(ctx, `
		SELECT achievement_id, unlocked_at
		FROM user_achievements
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]time.Time)

	for rows.Next() {

		var id int64
		var unlockedAt time.Time

		if err := rows.Scan(&id, &unlockedAt); err != nil {
			return nil, err
		}

		result[id] = unlockedAt
	}

	return result, rows.Err()
}

func (r *Repository) UnlockAchievement(ctx context.Context, userID, achievementID int64) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_achievements (user_id, achievement_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`, userID, achievementID)

	return err
}

func (r *Repository) GetWordsLearnedCount(ctx context.Context, userID int64) (int, error) {

	var count int

	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM word_learning_progress
		WHERE user_id = $1 AND learned_at IS NOT NULL
	`, userID).Scan(&count)

	return count, err
}

func (r *Repository) GetLessonsCompletedCount(ctx context.Context, userID int64) (int, error) {

	var count int

	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_lesson_progress
		WHERE user_id = $1 AND status = 'completed'
	`, userID).Scan(&count)

	return count, err
}
