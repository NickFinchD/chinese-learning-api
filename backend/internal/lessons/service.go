package lessons

import "context"

type Service struct {
	repository   *Repository
	wordProvider WordProvider
	quizProvider QuizProvider
}

func NewService(
	repository *Repository,
	wordProvider WordProvider,
	quizProvider QuizProvider,
) *Service {
	return &Service{
		repository:   repository,
		wordProvider: wordProvider,
		quizProvider: quizProvider,
	}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*LessonResponse, error) {

	lesson, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	steps, err := s.repository.GetSteps(ctx, id)
	if err != nil {
		return nil, err
	}

	loadedSteps, err := s.loadSteps(ctx, steps)
	if err != nil {
		return nil, err
	}

	return &LessonResponse{
		ID:           lesson.ID,
		Title:        lesson.Title,
		Description:  lesson.Description,
		LessonNumber: lesson.LessonNumber,
		Steps:        loadedSteps,
	}, nil
}
