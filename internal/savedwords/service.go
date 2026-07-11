package savedwords

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Save(ctx context.Context, userID, wordID int64) error {
	return s.repository.Save(ctx, userID, wordID)
}
func (s *Service) List(ctx context.Context, userID int64) ([]words.Word, error) {
	return s.repository.List(ctx, userID)
}
func (s *Service) Delete(ctx context.Context, userID, wordID int64) error {
	return s.repository.Delete(ctx, userID, wordID)
}
