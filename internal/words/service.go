package words

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context, request ListRequest) ([]Word, error) {
	return s.repository.List(ctx, request)
}
func (s *Service) GetByID(ctx context.Context, id int64) (*Word, error) {
	return s.repository.GetByID(ctx, id)
}
