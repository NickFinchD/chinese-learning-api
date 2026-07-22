package quizzes

type CreateQuizRequest struct {
	Question string                `json:"question" binding:"required"`
	Options  []CreateOptionRequest `json:"options" binding:"required,min=2"`
}

type CreateOptionRequest struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}
