package gamification

import (
	"context"
	"testing"
	"time"
)

type fakeRepository struct {
	stats            *UserStats
	achievements     []Achievement
	unlocked         map[int64]time.Time
	wordsLearned     int
	lessonsCompleted int

	unlockedIDs []int64
	xpAdded     int
}

func (f *fakeRepository) GetStats(ctx context.Context, userID int64) (*UserStats, error) {
	return f.stats, nil
}

func (f *fakeRepository) Heartbeat(ctx context.Context, userID int64, maxGapSeconds int64) error {
	return nil
}

func (f *fakeRepository) AddXP(ctx context.Context, userID int64, amount int) error {
	f.xpAdded += amount
	f.stats.XP += amount
	return nil
}

func (f *fakeRepository) ListAchievements(ctx context.Context) ([]Achievement, error) {
	return f.achievements, nil
}

func (f *fakeRepository) ListUnlockedAchievementIDs(ctx context.Context, userID int64) (map[int64]time.Time, error) {
	return f.unlocked, nil
}

func (f *fakeRepository) UnlockAchievement(ctx context.Context, userID, achievementID int64) error {
	f.unlockedIDs = append(f.unlockedIDs, achievementID)
	f.unlocked[achievementID] = time.Now()
	return nil
}

func (f *fakeRepository) GetWordsLearnedCount(ctx context.Context, userID int64) (int, error) {
	return f.wordsLearned, nil
}

func (f *fakeRepository) GetLessonsCompletedCount(ctx context.Context, userID int64) (int, error) {
	return f.lessonsCompleted, nil
}

func TestCheckAchievements_UnlocksWhenThresholdReached(t *testing.T) {

	repo := &fakeRepository{
		stats:        &UserStats{UserID: 1},
		unlocked:     map[int64]time.Time{},
		wordsLearned: 10,
		achievements: []Achievement{
			{ID: 1, Metric: MetricWordsLearned, Threshold: 10, XPReward: 50},
		},
	}

	service := NewService(repo)

	achievements, err := service.GetAchievements(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.unlockedIDs) != 1 || repo.unlockedIDs[0] != 1 {
		t.Fatalf("expected achievement 1 to be unlocked, got %v", repo.unlockedIDs)
	}

	if repo.xpAdded != 50 {
		t.Fatalf("expected 50 xp to be granted, got %d", repo.xpAdded)
	}

	if !achievements[0].Unlocked {
		t.Fatal("expected the returned achievement to be marked as unlocked")
	}
}

func TestCheckAchievements_DoesNotUnlockBelowThreshold(t *testing.T) {

	repo := &fakeRepository{
		stats:        &UserStats{UserID: 1},
		unlocked:     map[int64]time.Time{},
		wordsLearned: 5,
		achievements: []Achievement{
			{ID: 1, Metric: MetricWordsLearned, Threshold: 10, XPReward: 50},
		},
	}

	service := NewService(repo)

	achievements, err := service.GetAchievements(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.unlockedIDs) != 0 {
		t.Fatalf("expected no achievement to unlock, got %v", repo.unlockedIDs)
	}

	if achievements[0].Unlocked {
		t.Fatal("expected the achievement to remain locked")
	}
}

func TestCheckAchievements_AlreadyUnlockedIsNotRewardedAgain(t *testing.T) {

	repo := &fakeRepository{
		stats:        &UserStats{UserID: 1, XP: 999},
		unlocked:     map[int64]time.Time{1: time.Now()},
		wordsLearned: 100,
		achievements: []Achievement{
			{ID: 1, Metric: MetricWordsLearned, Threshold: 10, XPReward: 50},
		},
	}

	service := NewService(repo)

	if _, err := service.GetAchievements(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.xpAdded != 0 {
		t.Fatalf("expected no additional xp for an already-unlocked achievement, got %d", repo.xpAdded)
	}
}

func TestGetProgress_ComputesLevelFromXP(t *testing.T) {

	repo := &fakeRepository{
		stats:        &UserStats{UserID: 1, XP: 150, TotalSecondsActive: 7200},
		unlocked:     map[int64]time.Time{},
		achievements: []Achievement{},
	}

	service := NewService(repo)

	progress, err := service.GetProgress(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Level != 2 {
		t.Fatalf("expected level 2 at 150 xp, got %d", progress.Level)
	}

	if progress.TotalSecondsActive != 7200 {
		t.Fatalf("expected total seconds active to pass through, got %d", progress.TotalSecondsActive)
	}
}
