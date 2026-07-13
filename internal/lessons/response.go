package lessons

type LessonResponse struct {
	ID           int64                `json:"id"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	LessonNumber int                  `json:"lesson_number"`
	Steps        []LessonStepResponse `json:"steps"`
}
