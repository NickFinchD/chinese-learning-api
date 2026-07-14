package progress

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) StartLesson(
	ctx context.Context,
	userID int64,
	lessonID int64,
) error {

	return s.repository.StartLesson(
		ctx,
		userID,
		lessonID,
	)
}

func (s *Service) GetProgress(
	ctx context.Context,
	userID int64,
	lessonID int64,
) (*ProgressResponse, error) {

	progress, err := s.repository.GetProgress(
		ctx,
		userID,
		lessonID,
	)

	if err != nil {
		return nil, err
	}

	return &ProgressResponse{
		Status:      progress.Status,
		CurrentStep: progress.CurrentStep,
		Score:       progress.Score,
	}, nil
}
func (s *Service) UpdateStep(
	ctx context.Context,
	userID int64,
	lessonID int64,
	currentStep int,
) error {

	return s.repository.UpdateStep(
		ctx,
		userID,
		lessonID,
		currentStep,
	)
}
func (s *Service) CompleteLesson(
	ctx context.Context,
	userID int64,
	lessonID int64,
	score int,
) error {

	return s.repository.CompleteLesson(
		ctx,
		userID,
		lessonID,
		score,
	)
}
