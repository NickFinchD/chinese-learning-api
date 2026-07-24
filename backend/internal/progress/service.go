package progress

import "context"

// LessonCompletionXP is the flat XP reward for completing a lesson for the
// first time. Retakes (see stores/lessons.ts restart on the frontend) don't
// grant it again — see CompleteLesson.
const LessonCompletionXP = 10

type repository interface {
	StartLesson(ctx context.Context, userID, lessonID int64) error
	GetProgress(ctx context.Context, userID, lessonID int64) (*UserLessonProgress, error)
	UpdateStep(ctx context.Context, userID, lessonID int64, currentStep int) error
	CompleteLesson(ctx context.Context, userID, lessonID int64, score int) (alreadyCompleted bool, err error)
	UpdateCourseProgress(ctx context.Context, userID, lessonID int64) error
}

type xpAwarder interface {
	AwardXP(ctx context.Context, userID int64, amount int) error
}

type Service struct {
	repository repository
	xpAwarder  xpAwarder
}

func NewService(repository repository, xpAwarder xpAwarder) *Service {
	return &Service{
		repository: repository,
		xpAwarder:  xpAwarder,
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
// CompleteLesson marks the lesson completed and returns the XP just
// awarded (0 on a retake — only the first completion pays out).
func (s *Service) CompleteLesson(
	ctx context.Context,
	userID int64,
	lessonID int64,
	score int,
) (int, error) {

	alreadyCompleted, err := s.repository.CompleteLesson(
		ctx,
		userID,
		lessonID,
		score,
	)

	if err != nil {
		return 0, err
	}

	if err := s.repository.UpdateCourseProgress(ctx, userID, lessonID); err != nil {
		return 0, err
	}

	if alreadyCompleted {
		return 0, nil
	}

	if err := s.xpAwarder.AwardXP(ctx, userID, LessonCompletionXP); err != nil {
		return 0, err
	}

	return LessonCompletionXP, nil
}
