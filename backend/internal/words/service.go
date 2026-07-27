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

func (s *Service) AdminList(ctx context.Context, request AdminListRequest) ([]Word, int, error) {
	return s.repository.AdminList(ctx, request)
}

func (s *Service) Create(ctx context.Context, req CreateWordRequest) (*Word, error) {
	return s.repository.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateWordRequest) (*Word, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
