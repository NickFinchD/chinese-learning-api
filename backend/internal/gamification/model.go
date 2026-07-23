package gamification

import "time"

type UserStats struct {
	UserID             int64      `db:"user_id"`
	XP                 int        `db:"xp"`
	TotalSecondsActive int64      `db:"total_seconds_active"`
	LastHeartbeatAt    *time.Time `db:"last_heartbeat_at"`
}

type Achievement struct {
	ID          int64  `db:"id"`
	Code        string `db:"code"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Tier        int16  `db:"tier"`
	Metric      string `db:"metric"`
	Threshold   int    `db:"threshold"`
	XPReward    int    `db:"xp_reward"`
}

// Metrics an achievement can be tied to.
const (
	MetricHoursActive      = "hours_active"
	MetricWordsLearned     = "words_learned"
	MetricLessonsCompleted = "lessons_completed"
)

// XPForLevel returns the cumulative XP required to reach a given level.
// The curve is a simple triangular progression, so each level requires
// progressively more XP than the last (level 2 needs 100, level 3 needs
// 200 more on top of that, and so on).
func XPForLevel(level int) int {
	n := level - 1
	return 100 * n * (n + 1) / 2
}

// LevelForXP returns the current level for a given amount of XP, along with
// the XP earned within the current level and the XP needed to reach the next one.
func LevelForXP(xp int) (level int, currentLevelXP int, xpForNextLevel int) {

	level = 1

	for XPForLevel(level+1) <= xp {
		level++
	}

	currentThreshold := XPForLevel(level)
	nextThreshold := XPForLevel(level + 1)

	return level, xp - currentThreshold, nextThreshold - currentThreshold
}
