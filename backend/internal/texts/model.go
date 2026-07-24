package texts

import "time"

type Text struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Hanzi       string    `db:"hanzi" json:"hanzi"`
	Pinyin      string    `db:"pinyin" json:"pinyin"`
	Translation string    `db:"translation" json:"translation"`
	HSKLevel    int16     `db:"hsk_level" json:"hsk_level"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
