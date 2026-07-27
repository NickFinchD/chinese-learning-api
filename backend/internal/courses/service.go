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

func (s *Service) List(ctx context.Context, userID int64) ([]Course, error) {
	return s.repository.List(ctx, userID)
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

func (s *Service) AdminList(ctx context.Context, request AdminListRequest) ([]Course, int, error) {
	return s.repository.AdminList(ctx, request)
}

// AdminGetByID returns the raw Course row (including sort_order, which
// CourseDetailsResponse omits) for the admin edit form.
func (s *Service) AdminGetByID(ctx context.Context, id int64) (*Course, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CourseRequest) (*Course, error) {
	return s.repository.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id int64, req CourseRequest) (*Course, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
