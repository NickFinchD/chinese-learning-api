package review

import "time"

type UserWordProgress struct {
	ID           int64      `db:"id"`
	UserID       int64      `db:"user_id"`
	WordID       int64      `db:"word_id"`
	ReviewCount  int        `db:"review_count"`
	CorrectCount int        `db:"correct_count"`
	WrongCount   int        `db:"wrong_count"`
	EaseFactor   float64    `db:"ease_factor"`
	LastReviewAt *time.Time `db:"last_review_at"`
	NextReviewAt time.Time  `db:"next_review_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
