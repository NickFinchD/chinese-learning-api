package texts

import "context"

type repository interface {
	List(ctx context.Context, hskLevel int16) ([]Text, error)
	GetByID(ctx context.Context, id int64) (*Text, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context, hskLevel int16) ([]Text, error) {
	return s.repository.List(ctx, hskLevel)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Text, error) {
	return s.repository.GetByID(ctx, id)
}
