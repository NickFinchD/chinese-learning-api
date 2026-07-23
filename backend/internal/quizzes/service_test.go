package quizzes

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	quiz    *Quiz
	quizzes []Quiz
	err     error

	checkAnswerResult bool
	checkAnswerErr    error

	lastCheckedQuizID   int64
	lastCheckedOptionID int64
}

func (f *fakeRepository) GetByIDs(ctx context.Context, ids []int64) ([]Quiz, error) {
	return f.quizzes, f.err
}

func (f *fakeRepository) GetAll(ctx context.Context) ([]Quiz, error) {
	return f.quizzes, f.err
}

func (f *fakeRepository) GetByHSKLevel(ctx context.Context, hsk int16) ([]Quiz, error) {
	return f.quizzes, f.err
}

func (f *fakeRepository) GetByID(ctx context.Context, id int64) (*Quiz, error) {
	return f.quiz, f.err
}

func (f *fakeRepository) Create(ctx context.Context, quiz Quiz) (*Quiz, error) {
	return &quiz, f.err
}

func (f *fakeRepository) CheckAnswer(ctx context.Context, quizID, optionID int64) (bool, error) {
	f.lastCheckedQuizID = quizID
	f.lastCheckedOptionID = optionID

	return f.checkAnswerResult, f.checkAnswerErr
}

func TestCheckAnswer_ReturnsRepositoryResult(t *testing.T) {

	repo := &fakeRepository{checkAnswerResult: true}

	service := NewService(repo)

	correct, err := service.CheckAnswer(context.Background(), 10, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !correct {
		t.Fatal("expected the answer to be reported as correct")
	}

	if repo.lastCheckedQuizID != 10 || repo.lastCheckedOptionID != 20 {
		t.Fatalf("expected repository to receive quiz id 10 and option id 20, got %d/%d",
			repo.lastCheckedQuizID, repo.lastCheckedOptionID)
	}
}

func TestCheckAnswer_PropagatesError(t *testing.T) {

	repo := &fakeRepository{checkAnswerErr: errors.New("not found")}

	service := NewService(repo)

	_, err := service.CheckAnswer(context.Background(), 10, 20)

	if err == nil {
		t.Fatal("expected an error to propagate")
	}
}

func TestGetByID_PropagatesError(t *testing.T) {

	repo := &fakeRepository{err: errors.New("db error")}

	service := NewService(repo)

	_, err := service.GetByID(context.Background(), 1)

	if err == nil {
		t.Fatal("expected an error to propagate")
	}
}
