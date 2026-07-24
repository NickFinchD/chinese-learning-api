package names

import "time"

// Name is a proper noun (person's name, book title, place name) that
// appears in a reading text. Kept separate from words.Word on purpose: it
// gets a hover hint in the reader like any dictionary word, but isn't real
// vocabulary to learn/train/quiz on.
type Name struct {
	ID          int64     `db:"id" json:"id"`
	Hanzi       string    `db:"hanzi" json:"hanzi"`
	Pinyin      string    `db:"pinyin" json:"pinyin"`
	Translation string    `db:"translation" json:"translation"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
