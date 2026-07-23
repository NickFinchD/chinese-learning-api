package progress

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	getProgressResult *UserLessonProgress
	getProgressErr    error

	completeLessonErr        error
	completeLessonCalled     bool
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

func (f *fakeRepository) CompleteLesson(ctx context.Context, userID, lessonID int64, score int) error {
	f.completeLessonCalled = true
	return f.completeLessonErr
}

func (f *fakeRepository) UpdateCourseProgress(ctx context.Context, userID, lessonID int64) error {
	f.updateCourseProgressCalled = true
	return f.updateCourseProgressErr
}

func TestGetProgress_MapsRepositoryResult(t *testing.T) {

	repo := &fakeRepository{
		getProgressResult: &UserLessonProgress{
			Status:      "in_progress",
			CurrentStep: 3,
			Score:       0,
		},
	}

	service := NewService(repo)

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

	service := NewService(repo)

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

	service := NewService(repo)

	_, err := service.GetProgress(context.Background(), 1, 1)

	if err == nil {
		t.Fatal("expected an error to propagate")
	}
}

func TestCompleteLesson_UpdatesCourseProgressAfterSuccess(t *testing.T) {

	repo := &fakeRepository{}

	service := NewService(repo)

	err := service.CompleteLesson(context.Background(), 1, 1, 100)

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

	service := NewService(repo)

	err := service.CompleteLesson(context.Background(), 1, 1, 100)

	if err == nil {
		t.Fatal("expected an error")
	}

	if repo.updateCourseProgressCalled {
		t.Fatal("UpdateCourseProgress should not run when CompleteLesson fails")
	}
}
