package grammar

import "time"

type Note struct {
	ID                  int64     `db:"id" json:"id"`
	Title               string    `db:"title" json:"title"`
	Explanation         string    `db:"explanation" json:"explanation"`
	ExampleHanzi        string    `db:"example_hanzi" json:"example_hanzi,omitempty"`
	ExamplePinyin       string    `db:"example_pinyin" json:"example_pinyin,omitempty"`
	ExampleTranslation  string    `db:"example_translation" json:"example_translation,omitempty"`
	Example2Hanzi       string    `db:"example2_hanzi" json:"example2_hanzi,omitempty"`
	Example2Pinyin      string    `db:"example2_pinyin" json:"example2_pinyin,omitempty"`
	Example2Translation string    `db:"example2_translation" json:"example2_translation,omitempty"`
	Example3Hanzi       string    `db:"example3_hanzi" json:"example3_hanzi,omitempty"`
	Example3Pinyin      string    `db:"example3_pinyin" json:"example3_pinyin,omitempty"`
	Example3Translation string    `db:"example3_translation" json:"example3_translation,omitempty"`
	HSKLevel            int16     `db:"hsk_level" json:"hsk_level"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}
