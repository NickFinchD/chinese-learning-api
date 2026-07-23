package learning

import (
	"context"
	"testing"
	"time"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type fakeRepository struct {
	stored map[int64]*WordLearningProgress
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{stored: make(map[int64]*WordLearningProgress)}
}

func (f *fakeRepository) Get(ctx context.Context, userID, wordID int64) (*WordLearningProgress, error) {
	if p, ok := f.stored[wordID]; ok {
		copy := *p
		return &copy, nil
	}

	return nil, nil
}

func (f *fakeRepository) Upsert(ctx context.Context, progress *WordLearningProgress) error {
	copy := *progress
	f.stored[progress.WordID] = &copy

	return nil
}

func (f *fakeRepository) ListForUser(ctx context.Context, userID int64) ([]WordLearningProgress, error) {
	items := make([]WordLearningProgress, 0, len(f.stored))

	for _, p := range f.stored {
		items = append(items, *p)
	}

	return items, nil
}

func (f *fakeRepository) ListLearnedWords(ctx context.Context, userID int64) ([]words.Word, error) {
	return nil, nil
}

func (f *fakeRepository) ListInProgress(ctx context.Context, userID int64) ([]WordProgressDetail, error) {
	return nil, nil
}

func newTestService(repo *fakeRepository, now time.Time) *Service {
	service := NewService(repo)
	service.now = func() time.Time { return now }

	return service
}

func TestRecordAnswer_FirstCorrectAnswerStartsAtStageOne(t *testing.T) {

	repo := newFakeRepository()
	now := time.Now()
	service := newTestService(repo, now)

	progress, err := service.RecordAnswer(context.Background(), 1, 100, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Stage != 1 {
		t.Fatalf("expected stage 1, got %d", progress.Stage)
	}

	wantNext := now.Add(5 * time.Minute)

	if progress.NextEligibleAt == nil || !progress.NextEligibleAt.Equal(wantNext) {
		t.Fatalf("expected next eligible at %v, got %v", wantNext, progress.NextEligibleAt)
	}

	if progress.LearnedAt != nil {
		t.Fatal("word should not be learned yet")
	}
}

func TestRecordAnswer_TooSoonDoesNotAdvance(t *testing.T) {

	repo := newFakeRepository()
	now := time.Now()
	service := newTestService(repo, now)

	if _, err := service.RecordAnswer(context.Background(), 1, 100, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Answer again immediately, before the 5 minute cooldown elapses.
	progress, err := service.RecordAnswer(context.Background(), 1, 100, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Stage != 1 {
		t.Fatalf("expected stage to stay at 1, got %d", progress.Stage)
	}
}

func TestRecordAnswer_AdvancesAfterCooldownElapses(t *testing.T) {

	repo := newFakeRepository()
	now := time.Now()
	service := newTestService(repo, now)

	if _, err := service.RecordAnswer(context.Background(), 1, 100, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	service.now = func() time.Time { return now.Add(6 * time.Minute) }

	progress, err := service.RecordAnswer(context.Background(), 1, 100, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Stage != 2 {
		t.Fatalf("expected stage 2, got %d", progress.Stage)
	}
}

func TestRecordAnswer_IncorrectAnswerDoesNotAdvance(t *testing.T) {

	repo := newFakeRepository()
	now := time.Now()
	service := newTestService(repo, now)

	if _, err := service.RecordAnswer(context.Background(), 1, 100, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	service.now = func() time.Time { return now.Add(time.Hour) }

	progress, err := service.RecordAnswer(context.Background(), 1, 100, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Stage != 1 {
		t.Fatalf("expected stage to stay at 1 after a wrong answer, got %d", progress.Stage)
	}
}

func TestRecordAnswer_GraduatesToLearnedAfterFullSchedule(t *testing.T) {

	repo := newFakeRepository()
	now := time.Now()
	service := newTestService(repo, now)

	elapsed := now

	for stage := int16(1); stage <= MaxStage; stage++ {

		service.now = func() time.Time { return elapsed }

		progress, err := service.RecordAnswer(context.Background(), 1, 100, true)

		if err != nil {
			t.Fatalf("unexpected error at stage %d: %v", stage, err)
		}

		if progress.Stage != stage {
			t.Fatalf("expected stage %d, got %d", stage, progress.Stage)
		}

		if interval, ok := StageIntervals[stage]; ok {
			elapsed = elapsed.Add(interval + time.Minute)
		}
	}

	final, err := repo.Get(context.Background(), 1, 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if final.LearnedAt == nil {
		t.Fatal("expected the word to be marked as learned after completing the schedule")
	}

	if final.NextEligibleAt != nil {
		t.Fatal("a learned word should not have a pending next_eligible_at")
	}
}

func TestRecordAnswer_AlreadyLearnedWordIsUnaffected(t *testing.T) {

	repo := newFakeRepository()
	now := time.Now()
	learnedAt := now.Add(-time.Hour)

	repo.stored[100] = &WordLearningProgress{
		UserID:    1,
		WordID:    100,
		Stage:     MaxStage,
		LearnedAt: &learnedAt,
	}

	service := newTestService(repo, now)

	progress, err := service.RecordAnswer(context.Background(), 1, 100, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Stage != MaxStage {
		t.Fatalf("expected stage to remain %d, got %d", MaxStage, progress.Stage)
	}
}
