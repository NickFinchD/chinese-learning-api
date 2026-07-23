package quizzes

import "context"

type repository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]Quiz, error)
	GetAll(ctx context.Context) ([]Quiz, error)
	GetByHSKLevel(ctx context.Context, hsk int16) ([]Quiz, error)
	GetByID(ctx context.Context, id int64) (*Quiz, error)
	Create(ctx context.Context, quiz Quiz) (*Quiz, error)
	CheckAnswer(ctx context.Context, quizID, optionID int64) (bool, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByIDs(ctx context.Context, ids []int64) ([]Quiz, error) {
	return s.repository.GetByIDs(ctx, ids)
}

func (s *Service) GetAll(ctx context.Context) ([]Quiz, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) GetByHSKLevel(ctx context.Context, hsk int16) ([]Quiz, error) {
	return s.repository.GetByHSKLevel(ctx, hsk)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Quiz, error) {
	return s.repository.GetByID(ctx, id)
}
func (s *Service) Create(ctx context.Context, quiz Quiz) (*Quiz, error) {
	return s.repository.Create(ctx, quiz)
}
func (s *Service) CheckAnswer(ctx context.Context, quizID, optionID int64) (bool, error) {
	return s.repository.CheckAnswer(ctx, quizID, optionID)
}
