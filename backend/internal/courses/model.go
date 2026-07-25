package courses

import "time"

type Course struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	HSKLevel    int16     `db:"hsk_level" json:"hsk_level"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	// ProgressPercent is only populated by List (it's user-scoped); 0 on a
	// course the user hasn't started.
	ProgressPercent int `db:"progress_percent" json:"progress_percent"`
}
