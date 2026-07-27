package lessons

import "context"

type Service struct {
	repository       *Repository
	wordProvider     WordProvider
	quizProvider     QuizProvider
	grammarProvider  GrammarProvider
	sentenceProvider SentenceProvider
}

func NewService(
	repository *Repository,
	wordProvider WordProvider,
	quizProvider QuizProvider,
	grammarProvider GrammarProvider,
	sentenceProvider SentenceProvider,
) *Service {
	return &Service{
		repository:       repository,
		wordProvider:     wordProvider,
		quizProvider:     quizProvider,
		grammarProvider:  grammarProvider,
		sentenceProvider: sentenceProvider,
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
		CourseID:     lesson.CourseID,
		Title:        lesson.Title,
		Description:  lesson.Description,
		LessonNumber: lesson.LessonNumber,
		Steps:        loadedSteps,
	}, nil
}

func (s *Service) AdminListByCourse(ctx context.Context, courseID int64) ([]AdminLessonListItem, error) {
	return s.repository.AdminListByCourse(ctx, courseID)
}

func (s *Service) Create(ctx context.Context, req LessonRequest) (*Lesson, error) {
	return s.repository.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id int64, req LessonRequest) (*Lesson, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) AddStep(ctx context.Context, lessonID int64, stepType string, entityID int64) (*LessonStep, error) {
	return s.repository.AddStep(ctx, lessonID, stepType, entityID)
}

func (s *Service) RemoveStep(ctx context.Context, stepID int64) error {
	return s.repository.RemoveStep(ctx, stepID)
}

func (s *Service) ReorderSteps(ctx context.Context, lessonID int64, stepIDs []int64) error {
	return s.repository.ReorderSteps(ctx, lessonID, stepIDs)
}
