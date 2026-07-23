package quizzes

import "time"

type Quiz struct {
	ID        int64     `db:"id" json:"id"`
	Question  string    `db:"question" json:"question"`
	HSKLevel  int16     `db:"hsk_level" json:"hsk_level"`
	Options   []Option  `json:"options"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Option struct {
	ID        int64  `db:"id" json:"id"`
	QuizID    int64  `db:"quiz_id" json:"-"`
	Text      string `db:"option_text" json:"text"`
	IsCorrect bool   `db:"is_correct" json:"-"`
	SortOrder int16  `db:"sort_order" json:"-"`
}
