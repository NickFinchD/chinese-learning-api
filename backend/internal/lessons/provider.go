package lessons

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/grammar"
	"github.com/NickFinchD/chinese-learning-api/internal/quizzes"
	"github.com/NickFinchD/chinese-learning-api/internal/sentences"
	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type WordProvider interface {
	GetByIDs(ctx context.Context, ids []int64) ([]words.Word, error)
}
type QuizProvider interface {
	GetByIDs(ctx context.Context, ids []int64) ([]quizzes.Quiz, error)
}
type GrammarProvider interface {
	GetByIDs(ctx context.Context, ids []int64) ([]grammar.Note, error)
}
type SentenceProvider interface {
	GetByIDs(ctx context.Context, ids []int64) ([]sentences.Exercise, error)
}
