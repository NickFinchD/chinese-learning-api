package progress

import "time"

type UserCourseProgress struct {
	ID               int64      `db:"id" json:"id"`
	UserID           int64      `db:"user_id" json:"user_id"`
	CourseID         int64      `db:"course_id" json:"course_id"`
	CompletedLessons int        `db:"completed_lessons" json:"completed_lessons"`
	TotalLessons     int        `db:"total_lessons" json:"total_lessons"`
	ProgressPercent  int        `db:"progress_percent" json:"progress_percent"`
	CompletedAt      *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}
