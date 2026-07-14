package lessons

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type WordProvider interface {
	GetByIDs(ctx context.Context, ids []int64) ([]words.Word, error)
}
