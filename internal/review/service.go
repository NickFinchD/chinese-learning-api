package review

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetWordsForReview(
	ctx context.Context,
	userID int64,
) ([]ReviewWordResponse, error) {

	return s.repository.GetWordsForReview(
		ctx,
		userID,
	)
}
