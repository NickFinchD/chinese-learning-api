package collections

import "time"

type Collection struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"-"`
	Name      string    `db:"name" json:"name"`
	WordCount int       `db:"word_count" json:"word_count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
