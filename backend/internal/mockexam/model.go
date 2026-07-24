package mockexam

import (
	"time"

	"github.com/NickFinchD/chinese-learning-api/internal/quizzes"
	"github.com/NickFinchD/chinese-learning-api/internal/sentences"
)

// PassThresholdPercent mirrors the ~60% passing bar used by the real HSK
// exams, applied here to our simplified practice paper.
const PassThresholdPercent = 60

// ExamPassXP is awarded once per HSK level, the first time a user passes
// its mock exam — bigger than LessonCompletionXP since clearing a full
// timed paper is a bigger milestone than one lesson.
const ExamPassXP = 50

type Attempt struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"-"`
	HSKLevel        int16     `db:"hsk_level" json:"hsk_level"`
	TotalQuestions  int16     `db:"total_questions" json:"total_questions"`
	CorrectCount    int16     `db:"correct_count" json:"correct_count"`
	ScorePercent    int16     `db:"score_percent" json:"score_percent"`
	Passed          bool      `db:"passed" json:"passed"`
	DurationSeconds int32     `db:"duration_seconds" json:"duration_seconds"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// Question is one item on the generated paper: exactly one of Quiz or
// Sentence is set, per Type.
type Question struct {
	Type     string              `json:"type"`
	Quiz     *quizzes.Quiz       `json:"quiz,omitempty"`
	Sentence *sentences.Exercise `json:"sentence,omitempty"`
}

type Paper struct {
	HSKLevel         int16      `json:"hsk_level"`
	TimeLimitSeconds int        `json:"time_limit_seconds"`
	Questions        []Question `json:"questions"`
}

type QuizAnswer struct {
	QuizID   int64 `json:"quiz_id"`
	OptionID int64 `json:"option_id"`
}

type SentenceAnswer struct {
	ExerciseID int64    `json:"exercise_id"`
	Chunks     []string `json:"chunks"`
}

type Result struct {
	Attempt   Attempt `json:"attempt"`
	XPAwarded int     `json:"xp_awarded"`
}
