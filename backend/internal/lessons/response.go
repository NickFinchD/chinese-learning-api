package lessons

type LessonResponse struct {
	ID           int64                `json:"id"`
	CourseID     int64                `json:"course_id"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	LessonNumber int                  `json:"lesson_number"`
	Steps        []LessonStepResponse `json:"steps"`
}
