package lessons

type AdminLessonListItem struct {
	ID           int64  `json:"id"`
	CourseID     int64  `json:"course_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	LessonNumber int    `json:"lesson_number"`
	StepCount    int    `json:"step_count"`
}

// AdminStepResponse gives the admin step editor a human-readable label for
// each polymorphic step (the entity's hanzi/question/title/translation)
// instead of the full player-facing payload LessonStepResponse carries.
type AdminStepResponse struct {
	ID        int64  `json:"id"`
	StepType  string `json:"step_type"`
	EntityID  int64  `json:"entity_id"`
	SortOrder int    `json:"sort_order"`
	Label     string `json:"label"`
}
