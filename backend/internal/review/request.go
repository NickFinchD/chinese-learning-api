package review

type AddWordRequest struct {
	WordID int64 `json:"word_id" binding:"required"`
}
type AnswerRequest struct {
	WordID  int64 `json:"word_id" binding:"required"`
	Correct bool  `json:"correct"`
}
