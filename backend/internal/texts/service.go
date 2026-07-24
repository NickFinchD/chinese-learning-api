package texts

import "context"

type repository interface {
	List(ctx context.Context, hskLevel int16, userID int64) ([]Text, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Text, error)
	MarkStarted(ctx context.Context, userID, textID int64) error
	MarkRead(ctx context.Context, userID, textID int64) error
	MarkUnread(ctx context.Context, userID, textID int64) error
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context, hskLevel int16, userID int64) ([]Text, error) {
	return s.repository.List(ctx, hskLevel, userID)
}

// GetByID marks the text as opened (in_progress) before returning it, so
// the reader's list of "в процессе" texts reflects what they've viewed.
func (s *Service) GetByID(ctx context.Context, id int64, userID int64) (*Text, error) {

	if err := s.repository.MarkStarted(ctx, userID, id); err != nil {
		return nil, err
	}

	return s.repository.GetByID(ctx, id, userID)
}

func (s *Service) MarkRead(ctx context.Context, userID, textID int64) error {
	return s.repository.MarkRead(ctx, userID, textID)
}

func (s *Service) MarkUnread(ctx context.Context, userID, textID int64) error {
	return s.repository.MarkUnread(ctx, userID, textID)
}
