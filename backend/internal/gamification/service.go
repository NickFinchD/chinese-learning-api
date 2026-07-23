package gamification

import (
	"context"
	"time"
)

// maxHeartbeatGapSeconds caps how much time a single heartbeat can add,
// so a stale tab or a machine that went to sleep can't inflate the total.
const maxHeartbeatGapSeconds = 90

type repository interface {
	GetStats(ctx context.Context, userID int64) (*UserStats, error)
	Heartbeat(ctx context.Context, userID int64, maxGapSeconds int64) error
	AddXP(ctx context.Context, userID int64, amount int) error
	ListAchievements(ctx context.Context) ([]Achievement, error)
	ListUnlockedAchievementIDs(ctx context.Context, userID int64) (map[int64]time.Time, error)
	UnlockAchievement(ctx context.Context, userID, achievementID int64) error
	GetWordsLearnedCount(ctx context.Context, userID int64) (int, error)
	GetLessonsCompletedCount(ctx context.Context, userID int64) (int, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Heartbeat(ctx context.Context, userID int64) error {

	if err := s.repository.Heartbeat(ctx, userID, maxHeartbeatGapSeconds); err != nil {
		return err
	}

	return s.checkAchievements(ctx, userID)
}

type Progress struct {
	XP                 int
	Level              int
	CurrentLevelXP     int
	XPForNextLevel     int
	TotalSecondsActive int64
}

func (s *Service) GetProgress(ctx context.Context, userID int64) (*Progress, error) {

	if err := s.checkAchievements(ctx, userID); err != nil {
		return nil, err
	}

	stats, err := s.repository.GetStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	level, currentLevelXP, xpForNextLevel := LevelForXP(stats.XP)

	return &Progress{
		XP:                 stats.XP,
		Level:              level,
		CurrentLevelXP:     currentLevelXP,
		XPForNextLevel:     xpForNextLevel,
		TotalSecondsActive: stats.TotalSecondsActive,
	}, nil
}

type AchievementStatus struct {
	Code        string
	Title       string
	Description string
	Tier        int16
	Metric      string
	Threshold   int
	XPReward    int
	Unlocked    bool
	UnlockedAt  *time.Time
}

func (s *Service) GetAchievements(ctx context.Context, userID int64) ([]AchievementStatus, error) {

	if err := s.checkAchievements(ctx, userID); err != nil {
		return nil, err
	}

	catalog, err := s.repository.ListAchievements(ctx)
	if err != nil {
		return nil, err
	}

	unlocked, err := s.repository.ListUnlockedAchievementIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]AchievementStatus, 0, len(catalog))

	for _, a := range catalog {

		unlockedAt, isUnlocked := unlocked[a.ID]

		status := AchievementStatus{
			Code:        a.Code,
			Title:       a.Title,
			Description: a.Description,
			Tier:        a.Tier,
			Metric:      a.Metric,
			Threshold:   a.Threshold,
			XPReward:    a.XPReward,
			Unlocked:    isUnlocked,
		}

		if isUnlocked {
			t := unlockedAt
			status.UnlockedAt = &t
		}

		result = append(result, status)
	}

	return result, nil
}

// checkAchievements recomputes the user's current metrics and unlocks any
// achievement whose threshold has now been reached, granting its XP reward.
// It's safe to call often: already-unlocked achievements are skipped via
// their unique (user_id, achievement_id) constraint.
func (s *Service) checkAchievements(ctx context.Context, userID int64) error {

	catalog, err := s.repository.ListAchievements(ctx)
	if err != nil {
		return err
	}

	unlocked, err := s.repository.ListUnlockedAchievementIDs(ctx, userID)
	if err != nil {
		return err
	}

	metrics, err := s.currentMetrics(ctx, userID)
	if err != nil {
		return err
	}

	for _, a := range catalog {

		if _, done := unlocked[a.ID]; done {
			continue
		}

		value, known := metrics[a.Metric]
		if !known || value < a.Threshold {
			continue
		}

		if err := s.repository.UnlockAchievement(ctx, userID, a.ID); err != nil {
			return err
		}

		if err := s.repository.AddXP(ctx, userID, a.XPReward); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) currentMetrics(ctx context.Context, userID int64) (map[string]int, error) {

	stats, err := s.repository.GetStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	wordsLearned, err := s.repository.GetWordsLearnedCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	lessonsCompleted, err := s.repository.GetLessonsCompletedCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		MetricHoursActive:      int(stats.TotalSecondsActive / 3600),
		MetricWordsLearned:     wordsLearned,
		MetricLessonsCompleted: lessonsCompleted,
	}, nil
}
