package courses

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context) ([]Course, error) {
	return s.repository.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*CourseDetailsResponse, error) {

	course, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	lessons, err := s.repository.GetLessons(ctx, id)
	if err != nil {
		return nil, err
	}

	return &CourseDetailsResponse{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		HSKLevel:    course.HSKLevel,
		Lessons:     lessons,
	}, nil
}
