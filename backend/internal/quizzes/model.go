package quizzes

import "time"

const (
	DirectionWordToTranslation = "word_to_translation"
	DirectionTranslationToWord = "translation_to_word"
)

type Quiz struct {
	ID        int64  `db:"id" json:"id"`
	Question  string `db:"question" json:"question"`
	HSKLevel  int16  `db:"hsk_level" json:"hsk_level"`
	Direction string `db:"direction" json:"direction"`
	// Hanzi/Pinyin of the word being tested, when the quiz was generated from
	// one (either quiz direction). Nil for hand-written quizzes.
	Hanzi     *string   `db:"hanzi" json:"hanzi,omitempty"`
	Pinyin    *string   `db:"pinyin" json:"pinyin,omitempty"`
	Options   []Option  `json:"options"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Option struct {
	ID     int64  `db:"id" json:"id"`
	QuizID int64  `db:"quiz_id" json:"-"`
	Text   string `db:"option_text" json:"text"`
	// Pinyin of Text, populated only for translation_to_word options (Text
	// is a hanzi word there, as opposed to a translation).
	Pinyin    *string `db:"pinyin" json:"pinyin,omitempty"`
	IsCorrect bool    `db:"is_correct" json:"-"`
	SortOrder int16   `db:"sort_order" json:"-"`
}
