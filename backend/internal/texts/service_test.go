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
}

func (f *fakeRepository) List(ctx context.Context, hskLevel int16) ([]Text, error) {
	f.lastHSKLevel = hskLevel
	return f.texts, f.err
}

func (f *fakeRepository) GetByID(ctx context.Context, id int64) (*Text, error) {
	return f.text, f.err
}

func TestList_PassesHSKLevelThrough(t *testing.T) {

	repo := &fakeRepository{texts: []Text{{ID: 1, Title: "Test"}}}

	service := NewService(repo)

	result, err := service.List(context.Background(), 3)

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

func TestGetByID_PropagatesError(t *testing.T) {

	repo := &fakeRepository{err: errors.New("not found")}

	service := NewService(repo)

	_, err := service.GetByID(context.Background(), 1)

	if err == nil {
		t.Fatal("expected an error to propagate")
	}
}
