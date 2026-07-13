package progress

type UpdateStepRequest struct {
	CurrentStep int `json:"current_step" binding:"required,min=1"`
}
type CompleteLessonRequest struct {
	Score int `json:"score" binding:"min=0"`
}
