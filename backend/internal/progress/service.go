package progress

import "context"

type repository interface {
	StartLesson(ctx context.Context, userID, lessonID int64) error
	GetProgress(ctx context.Context, userID, lessonID int64) (*UserLessonProgress, error)
	UpdateStep(ctx context.Context, userID, lessonID int64, currentStep int) error
	CompleteLesson(ctx context.Context, userID, lessonID int64, score int) error
	UpdateCourseProgress(ctx context.Context, userID, lessonID int64) error
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
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

	if progress == nil {
		return &ProgressResponse{
			Status:      "not_started",
			CurrentStep: 0,
			Score:       0,
		}, nil
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

	err := s.repository.CompleteLesson(
		ctx,
		userID,
		lessonID,
		score,
	)

	if err != nil {
		return err
	}

	return s.repository.UpdateCourseProgress(
		ctx,
		userID,
		lessonID,
	)
}
