package learning

import "time"

const MaxStage int16 = 8

// StageIntervals[n] is how long to wait after reaching stage n
// before the next correct answer is allowed to advance to stage n+1.
var StageIntervals = map[int16]time.Duration{
	1: 5 * time.Minute,
	2: time.Hour,
	3: 6 * time.Hour,
	4: 24 * time.Hour,
	5: 48 * time.Hour,
	6: 96 * time.Hour,
	7: 7 * 24 * time.Hour,
}

type WordLearningProgress struct {
	ID             int64      `db:"id"`
	UserID         int64      `db:"user_id"`
	WordID         int64      `db:"word_id"`
	Stage          int16      `db:"stage"`
	NextEligibleAt *time.Time `db:"next_eligible_at"`
	LearnedAt      *time.Time `db:"learned_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// WordProgressDetail is a word joined with its in-progress learning state,
// used for the "in progress" list (words that have started the repetition
// schedule but haven't graduated to learned yet).
type WordProgressDetail struct {
	WordID         int64
	Hanzi          string
	Pinyin         string
	Translation    string
	HSKLevel       int16
	Stage          int16
	NextEligibleAt *time.Time
}
