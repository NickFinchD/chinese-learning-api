package sentences

import "time"

type Exercise struct {
	ID          int64     `db:"id" json:"id"`
	Translation string    `db:"translation" json:"translation"`
	Chunks      []string  `db:"chunks" json:"chunks"`
	Pinyin      string    `db:"pinyin" json:"pinyin"`
	HSKLevel    int16     `db:"hsk_level" json:"hsk_level"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
