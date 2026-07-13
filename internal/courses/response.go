package courses

type CourseDetailsResponse struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	HSKLevel    int16       `json:"hsk_level"`
	Lessons     []LessonDTO `json:"lessons"`
}

type LessonDTO struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	LessonNumber int    `json:"lesson_number"`
}
