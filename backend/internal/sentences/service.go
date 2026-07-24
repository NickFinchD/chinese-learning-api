package sentences

import "context"

type repository interface {
	List(ctx context.Context, hskLevel int16) ([]Exercise, error)
	GetByID(ctx context.Context, id int64) (*Exercise, error)
	GetByIDs(ctx context.Context, ids []int64) ([]Exercise, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context, hskLevel int16) ([]Exercise, error) {
	return s.repository.List(ctx, hskLevel)
}

func (s *Service) GetByIDs(ctx context.Context, ids []int64) ([]Exercise, error) {
	return s.repository.GetByIDs(ctx, ids)
}

// CheckAnswer reports whether orderedChunks reconstructs the exercise's
// sentence exactly, in order.
func (s *Service) CheckAnswer(ctx context.Context, exerciseID int64, orderedChunks []string) (bool, error) {

	exercise, err := s.repository.GetByID(ctx, exerciseID)
	if err != nil {
		return false, err
	}

	if len(orderedChunks) != len(exercise.Chunks) {
		return false, nil
	}

	for i, chunk := range orderedChunks {
		if chunk != exercise.Chunks[i] {
			return false, nil
		}
	}

	return true, nil
}
