package learning

import (
	"context"
	"time"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type repository interface {
	Get(ctx context.Context, userID, wordID int64) (*WordLearningProgress, error)
	Upsert(ctx context.Context, progress *WordLearningProgress) error
	ListForUser(ctx context.Context, userID int64) ([]WordLearningProgress, error)
	ListInProgress(ctx context.Context, userID int64) ([]WordProgressDetail, error)
	ListLearnedWords(ctx context.Context, userID int64) ([]words.Word, error)
}

type Service struct {
	repository repository
	now        func() time.Time
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
		now:        time.Now,
	}
}

// RecordAnswer advances a word's learning progress when the user answers
// it correctly in the word-training test, following a fixed schedule:
// stage 1 (immediate) -> 5m -> stage 2 -> 1h -> stage 3 -> 6h -> stage 4 ->
// 24h -> stage 5 -> 48h -> stage 6 -> 96h -> stage 7 -> 1 week -> stage 8 (learned).
// A correct answer given before the word's cooldown elapses does not advance
// it, and an incorrect answer never rolls progress back.
func (s *Service) RecordAnswer(ctx context.Context, userID, wordID int64, correct bool) (*WordLearningProgress, error) {

	progress, err := s.repository.Get(ctx, userID, wordID)
	if err != nil {
		return nil, err
	}

	if progress == nil {
		progress = &WordLearningProgress{
			UserID: userID,
			WordID: wordID,
			Stage:  0,
		}
	}

	if !correct || progress.LearnedAt != nil {
		if progress.ID == 0 {
			if err := s.repository.Upsert(ctx, progress); err != nil {
				return nil, err
			}
		}

		return progress, nil
	}

	now := s.now()

	if progress.Stage > 0 && progress.NextEligibleAt != nil && now.Before(*progress.NextEligibleAt) {
		return progress, nil
	}

	progress.Stage++

	if progress.Stage >= MaxStage {
		progress.Stage = MaxStage
		progress.NextEligibleAt = nil
		learnedAt := now
		progress.LearnedAt = &learnedAt
	} else {
		next := now.Add(StageIntervals[progress.Stage])
		progress.NextEligibleAt = &next
	}

	if err := s.repository.Upsert(ctx, progress); err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *Service) ListForUser(ctx context.Context, userID int64) ([]WordLearningProgress, error) {
	return s.repository.ListForUser(ctx, userID)
}

func (s *Service) ListLearnedWords(ctx context.Context, userID int64) ([]words.Word, error) {
	return s.repository.ListLearnedWords(ctx, userID)
}

func (s *Service) ListInProgress(ctx context.Context, userID int64) ([]WordProgressDetail, error) {
	return s.repository.ListInProgress(ctx, userID)
}
