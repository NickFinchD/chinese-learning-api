package savedwords

import "time"

type SavedWord struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	WordID    int64     `db:"word_id" json:"word_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
