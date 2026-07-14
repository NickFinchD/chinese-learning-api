package review

import (
	"context"

	"github.com/NickFinchD/chinese-learning-api/internal/words"
)

type WordProvider interface {
	GetByID(
		ctx context.Context,
		id int64,
	) (*words.Word, error)
}
