package quizzes

type CheckAnswerRequest struct {
	OptionID int64 `json:"option_id" binding:"required"`
}
