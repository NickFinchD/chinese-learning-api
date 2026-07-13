package lessons

type LessonStepResponse struct {
	ID        int64  `json:"id"`
	StepType  string `json:"step_type"`
	SortOrder int    `json:"sort_order"`
	Data      any    `json:"data"`
}
