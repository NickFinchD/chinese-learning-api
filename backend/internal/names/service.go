package names

import "context"

type repository interface {
	GetAll(ctx context.Context) ([]Name, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Name, error) {
	return s.repository.GetAll(ctx)
}
