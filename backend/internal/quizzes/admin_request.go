package quizzes

type AdminListRequest struct {
	Search string `form:"search"`
	HSK    int16  `form:"hsk"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type AdminOptionRequest struct {
	Text      string `json:"text" binding:"required"`
	Pinyin    string `json:"pinyin"`
	IsCorrect bool   `json:"is_correct"`
}

type AdminQuizRequest struct {
	Question  string               `json:"question" binding:"required"`
	HSKLevel  int16                `json:"hsk_level" binding:"required"`
	Direction string               `json:"direction" binding:"required"`
	Hanzi     string               `json:"hanzi"`
	Pinyin    string               `json:"pinyin"`
	Options   []AdminOptionRequest `json:"options" binding:"required,min=2,dive"`
}
