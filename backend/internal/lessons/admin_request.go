package lessons

type LessonRequest struct {
	CourseID     int64  `json:"course_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	LessonNumber int    `json:"lesson_number" binding:"required"`
}

type AddStepRequest struct {
	StepType string `json:"step_type" binding:"required"`
	EntityID int64  `json:"entity_id" binding:"required"`
}

type ReorderStepsRequest struct {
	StepIDs []int64 `json:"step_ids" binding:"required"`
}
