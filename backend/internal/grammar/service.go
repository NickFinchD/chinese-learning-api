package grammar

import "context"

type repository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]Note, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByIDs(ctx context.Context, ids []int64) ([]Note, error) {
	return s.repository.GetByIDs(ctx, ids)
}
