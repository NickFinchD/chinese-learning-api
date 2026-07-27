package grammar

import "context"

type repository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]Note, error)
	AdminList(ctx context.Context, request AdminListRequest) ([]Note, int, error)
	Create(ctx context.Context, req CreateNoteRequest) (*Note, error)
	Update(ctx context.Context, id int64, req UpdateNoteRequest) (*Note, error)
	Delete(ctx context.Context, id int64) error
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

func (s *Service) AdminList(ctx context.Context, request AdminListRequest) ([]Note, int, error) {
	return s.repository.AdminList(ctx, request)
}

func (s *Service) Create(ctx context.Context, req CreateNoteRequest) (*Note, error) {
	return s.repository.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateNoteRequest) (*Note, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
