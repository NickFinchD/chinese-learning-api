package collections

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type repository interface {
	Create(ctx context.Context, userID int64, name string) (*Collection, error)
	List(ctx context.Context, userID int64) ([]Collection, error)
	ListCurated(ctx context.Context) ([]Collection, error)
	GetByID(ctx context.Context, userID, collectionID int64) (*Collection, error)
	ListWords(ctx context.Context, collectionID int64) ([]words.Word, error)
	Rename(ctx context.Context, userID, collectionID int64, name string) error
	Delete(ctx context.Context, userID, collectionID int64) error
	AddWord(ctx context.Context, userID, collectionID, wordID int64) error
	RemoveWord(ctx context.Context, userID, collectionID, wordID int64) error
	CloneForUser(ctx context.Context, userID, curatedID int64) (*Collection, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

// Detail is a collection together with the words currently in it.
type Detail struct {
	Collection
	Words []words.Word `json:"words"`
}

func (s *Service) Create(ctx context.Context, userID int64, name string) (*Collection, error) {
	return s.repository.Create(ctx, userID, name)
}

func (s *Service) List(ctx context.Context, userID int64) ([]Collection, error) {
	return s.repository.List(ctx, userID)
}

func (s *Service) ListCurated(ctx context.Context) ([]Collection, error) {
	return s.repository.ListCurated(ctx)
}

func (s *Service) SaveCurated(ctx context.Context, userID, curatedID int64) (*Collection, error) {
	return s.repository.CloneForUser(ctx, userID, curatedID)
}

// GetByID returns nil (no error) when the collection doesn't exist or isn't
// owned by userID, so the handler can turn that into a 404.
func (s *Service) GetByID(ctx context.Context, userID, collectionID int64) (*Detail, error) {

	c, err := s.repository.GetByID(ctx, userID, collectionID)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	w, err := s.repository.ListWords(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	return &Detail{Collection: *c, Words: w}, nil
}

func (s *Service) Rename(ctx context.Context, userID, collectionID int64, name string) error {
	return s.repository.Rename(ctx, userID, collectionID, name)
}

func (s *Service) Delete(ctx context.Context, userID, collectionID int64) error {
	return s.repository.Delete(ctx, userID, collectionID)
}

func (s *Service) AddWord(ctx context.Context, userID, collectionID, wordID int64) error {
	return s.repository.AddWord(ctx, userID, collectionID, wordID)
}

func (s *Service) RemoveWord(ctx context.Context, userID, collectionID, wordID int64) error {
	return s.repository.RemoveWord(ctx, userID, collectionID, wordID)
}
