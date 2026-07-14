package review

import (
	"context"
	"time"
)

type Service struct {
	repository   *Repository
	wordProvider WordProvider
}

func NewService(
	repository *Repository,
	wordProvider WordProvider,
) *Service {

	return &Service{
		repository:   repository,
		wordProvider: wordProvider,
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
func (s *Service) AddWord(
	ctx context.Context,
	userID int64,
	wordID int64,
) error {

	_, err := s.wordProvider.GetByID(
		ctx,
		wordID,
	)

	if err != nil {
		return err
	}

	return s.repository.AddWord(
		ctx,
		userID,
		wordID,
	)
}
func (s *Service) Answer(
	ctx context.Context,
	userID int64,
	wordID int64,
	correct bool,
) error {

	progress, err := s.repository.GetProgress(
		ctx,
		userID,
		wordID,
	)

	if err != nil {
		return err
	}
	if progress.NextReviewAt.After(time.Now()) {
		return ErrWordNotReady
	}
	progress.ReviewCount++

	if correct {
		progress.CorrectCount++
	} else {
		progress.WrongCount++
	}

	progress.NextReviewAt = calculateNextReview(
		progress,
		correct,
	)

	progress.LastReviewAt = func() *time.Time {
		t := time.Now()
		return &t
	}()

	return s.repository.UpdateProgress(
		ctx,
		userID,
		wordID,
		progress.ReviewCount,
		progress.CorrectCount,
		progress.WrongCount,
		progress.NextReviewAt,
		*progress.LastReviewAt,
	)
}
func (s *Service) GetStatistics(
	ctx context.Context,
	userID int64,
) (*StatisticsResponse, error) {

	return s.repository.GetStatistics(
		ctx,
		userID,
	)
}
func (s *Service) StartSession(
	ctx context.Context,
	userID int64,
) (*ReviewSessionResponse, error) {

	words, err := s.repository.GetWordsForReview(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return &ReviewSessionResponse{
		Total: len(words),
		Words: words,
	}, nil
}
