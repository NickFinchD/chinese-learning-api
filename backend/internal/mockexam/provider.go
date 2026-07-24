package mockexam

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/quizzes"
	"github.com/NickFinchD/chinese-learning-api/internal/sentences"
)

type QuizProvider interface {
	GetByHSKLevel(ctx context.Context, hsk int16) ([]quizzes.Quiz, error)
	CheckAnswer(ctx context.Context, quizID, optionID int64) (bool, error)
}

type SentenceProvider interface {
	List(ctx context.Context, hskLevel int16) ([]sentences.Exercise, error)
	CheckAnswer(ctx context.Context, exerciseID int64, orderedChunks []string) (bool, error)
}

type XPAwarder interface {
	AwardXP(ctx context.Context, userID int64, amount int) error
}
