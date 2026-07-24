package texts

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	texts []Text
	text  *Text
	err   error

	lastHSKLevel int16

	markStartedCalled bool
	markReadCalled    bool
	markUnreadCalled  bool
}

func (f *fakeRepository) List(ctx context.Context, hskLevel int16, userID int64) ([]Text, error) {
	f.lastHSKLevel = hskLevel
	return f.texts, f.err
}

func (f *fakeRepository) GetByID(ctx context.Context, id int64, userID int64) (*Text, error) {
	return f.text, f.err
}

func (f *fakeRepository) MarkStarted(ctx context.Context, userID, textID int64) error {
	f.markStartedCalled = true
	return f.err
}

func (f *fakeRepository) MarkRead(ctx context.Context, userID, textID int64) error {
	f.markReadCalled = true
	return f.err
}

func (f *fakeRepository) MarkUnread(ctx context.Context, userID, textID int64) error {
	f.markUnreadCalled = true
	return f.err
}

func TestList_PassesHSKLevelThrough(t *testing.T) {

	repo := &fakeRepository{texts: []Text{{ID: 1, Title: "Test"}}}

	service := NewService(repo)

	result, err := service.List(context.Background(), 3, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastHSKLevel != 3 {
		t.Fatalf("expected hsk level 3 to be passed through, got %d", repo.lastHSKLevel)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 text, got %d", len(result))
	}
}

func TestGetByID_MarksStartedThenFetches(t *testing.T) {

	repo := &fakeRepository{text: &Text{ID: 1, Title: "Test"}}

	service := NewService(repo)

	result, err := service.GetByID(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.markStartedCalled {
		t.Fatal("expected MarkStarted to be called before fetching")
	}

	if result.ID != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestGetByID_PropagatesError(t *testing.T) {

	repo := &fakeRepository{err: errors.New("not found")}

	service := NewService(repo)

	_, err := service.GetByID(context.Background(), 1, 1)

	if err == nil {
		t.Fatal("expected an error to propagate")
	}
}

func TestMarkRead_DelegatesToRepository(t *testing.T) {

	repo := &fakeRepository{}

	service := NewService(repo)

	if err := service.MarkRead(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.markReadCalled {
		t.Fatal("expected MarkRead to be called")
	}
}

func TestMarkUnread_DelegatesToRepository(t *testing.T) {

	repo := &fakeRepository{}

	service := NewService(repo)

	if err := service.MarkUnread(context.Background(), 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.markUnreadCalled {
		t.Fatal("expected MarkUnread to be called")
	}
}
