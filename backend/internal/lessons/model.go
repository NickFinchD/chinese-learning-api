package lessons

import "time"

type Lesson struct {
	ID           int64     `db:"id" json:"id"`
	CourseID     int64     `db:"course_id" json:"course_id"`
	Title        string    `db:"title" json:"title"`
	Description  string    `db:"description" json:"description"`
	LessonNumber int       `db:"lesson_number" json:"lesson_number"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type LessonStep struct {
	ID        int64     `db:"id" json:"id"`
	LessonID  int64     `db:"lesson_id" json:"lesson_id"`
	StepType  string    `db:"step_type" json:"step_type"`
	EntityID  *int64    `db:"entity_id" json:"entity_id,omitempty"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
