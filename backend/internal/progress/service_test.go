package progress

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	getProgressResult *UserLessonProgress
	getProgressErr    error

	completeLessonErr          error
	completeLessonAlreadyDone  bool
	completeLessonCalled       bool
	updateCourseProgressCalled bool
	updateCourseProgressErr    error
}

func (f *fakeRepository) StartLesson(ctx context.Context, userID, lessonID int64) error {
	return nil
}

func (f *fakeRepository) GetProgress(ctx context.Context, userID, lessonID int64) (*UserLessonProgress, error) {
	return f.getProgressResult, f.getProgressErr
}

func (f *fakeRepository) UpdateStep(ctx context.Context, userID, lessonID int64, currentStep int) error {
	return nil
}

func (f *fakeRepository) CompleteLesson(ctx context.Context, userID, lessonID int64, score int) (bool, error) {
	f.completeLessonCalled = true
	return f.completeLessonAlreadyDone, f.completeLessonErr
}

func (f *fakeRepository) UpdateCourseProgress(ctx context.Context, userID, lessonID int64) error {
	f.updateCourseProgressCalled = true
	return f.updateCourseProgressErr
}

type fakeXPAwarder struct {
	awardedAmount int
	callCount     int
	err           error
}

func (f *fakeXPAwarder) AwardXP(ctx context.Context, userID int64, amount int) error {
	f.callCount++
	f.awardedAmount = amount
	return f.err
}

func TestGetProgress_MapsRepositoryResult(t *testing.T) {

	repo := &fakeRepository{
		getProgressResult: &UserLessonProgress{
			Status:      "in_progress",
			CurrentStep: 3,
			Score:       0,
		},
	}

	service := NewService(repo, &fakeXPAwarder{})

	result, err := service.GetProgress(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "in_progress" || result.CurrentStep != 3 {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestGetProgress_ReturnsNotStartedWhenNoRowExists(t *testing.T) {

	repo := &fakeRepository{}

	service := NewService(repo, &fakeXPAwarder{})

	result, err := service.GetProgress(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "not_started" || result.CurrentStep != 0 {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestGetProgress_PropagatesError(t *testing.T) {

	repo := &fakeRepository{
		getProgressErr: errors.New("not found"),
	}

	service := NewService(repo, &fakeXPAwarder{})

	_, err := service.GetProgress(context.Background(), 1, 1)

	if err == nil {
		t.Fatal("expected an error to propagate")
	}
}

func TestCompleteLesson_UpdatesCourseProgressAfterSuccess(t *testing.T) {

	repo := &fakeRepository{}

	service := NewService(repo, &fakeXPAwarder{})

	_, err := service.CompleteLesson(context.Background(), 1, 1, 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.completeLessonCalled {
		t.Fatal("expected CompleteLesson to be called")
	}

	if !repo.updateCourseProgressCalled {
		t.Fatal("expected UpdateCourseProgress to be called after a successful completion")
	}
}

func TestCompleteLesson_SkipsCourseProgressOnFailure(t *testing.T) {

	repo := &fakeRepository{
		completeLessonErr: errors.New("db error"),
	}

	service := NewService(repo, &fakeXPAwarder{})

	_, err := service.CompleteLesson(context.Background(), 1, 1, 100)

	if err == nil {
		t.Fatal("expected an error")
	}

	if repo.updateCourseProgressCalled {
		t.Fatal("UpdateCourseProgress should not run when CompleteLesson fails")
	}
}

func TestCompleteLesson_AwardsXPOnFirstCompletion(t *testing.T) {

	repo := &fakeRepository{completeLessonAlreadyDone: false}
	xp := &fakeXPAwarder{}

	service := NewService(repo, xp)

	awarded, err := service.CompleteLesson(context.Background(), 1, 1, 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if xp.callCount != 1 || xp.awardedAmount != LessonCompletionXP {
		t.Fatalf("expected %d XP to be awarded once, got %d call(s) of %d", LessonCompletionXP, xp.callCount, xp.awardedAmount)
	}

	if awarded != LessonCompletionXP {
		t.Fatalf("expected CompleteLesson to return %d, got %d", LessonCompletionXP, awarded)
	}
}

func TestCompleteLesson_SkipsXPOnRetake(t *testing.T) {

	repo := &fakeRepository{completeLessonAlreadyDone: true}
	xp := &fakeXPAwarder{}

	service := NewService(repo, xp)

	awarded, err := service.CompleteLesson(context.Background(), 1, 1, 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if xp.callCount != 0 {
		t.Fatalf("expected no XP award on a repeat completion, got %d call(s)", xp.callCount)
	}

	if awarded != 0 {
		t.Fatalf("expected CompleteLesson to return 0, got %d", awarded)
	}
}
