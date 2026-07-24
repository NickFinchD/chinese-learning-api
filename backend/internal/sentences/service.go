package sentences

import "context"

type repository interface {
	List(ctx context.Context, hskLevel int16) ([]Exercise, error)
	GetByIDs(ctx context.Context, ids []int64) ([]Exercise, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context, hskLevel int16) ([]Exercise, error) {
	return s.repository.List(ctx, hskLevel)
}

func (s *Service) GetByIDs(ctx context.Context, ids []int64) ([]Exercise, error) {
	return s.repository.GetByIDs(ctx, ids)
}
