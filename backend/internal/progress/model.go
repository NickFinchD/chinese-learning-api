package progress

import "time"

type UserLessonProgress struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"user_id"`
	LessonID    int64      `db:"lesson_id" json:"lesson_id"`
	Status      string     `db:"status" json:"status"`
	CurrentStep int        `db:"current_step" json:"current_step"`
	Score       int        `db:"score" json:"score"`
	StartedAt   *time.Time `db:"started_at" json:"started_at,omitempty"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
