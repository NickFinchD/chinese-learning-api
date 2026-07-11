package words

import "time"

type Word struct {
	ID           int64     `db:"id" json:"id"`
	Hanzi        string    `db:"hanzi" json:"hanzi"`
	Pinyin       string    `db:"pinyin" json:"pinyin"`
	Translation  string    `db:"translation" json:"translation"`
	PartOfSpeech string    `db:"part_of_speech" json:"part_of_speech"`
	HSKLevel     int16     `db:"hsk_level" json:"hsk_level"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
