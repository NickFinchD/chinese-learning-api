package quizzes

type QuizResponse struct {
	ID       int64            `json:"id"`
	Question string           `json:"question"`
	Options  []OptionResponse `json:"options"`
}

type OptionResponse struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
	SortOrder int16  `json:"sort_order"`
}
