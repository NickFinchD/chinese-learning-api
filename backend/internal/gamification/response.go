package gamification

import "time"

type ProgressResponse struct {
	XP                 int   `json:"xp"`
	Level              int   `json:"level"`
	CurrentLevelXP     int   `json:"current_level_xp"`
	XPForNextLevel     int   `json:"xp_for_next_level"`
	TotalSecondsActive int64 `json:"total_seconds_active"`
	HoursActive        int   `json:"hours_active"`
}

func toProgressResponse(p *Progress) ProgressResponse {
	return ProgressResponse{
		XP:                 p.XP,
		Level:              p.Level,
		CurrentLevelXP:     p.CurrentLevelXP,
		XPForNextLevel:     p.XPForNextLevel,
		TotalSecondsActive: p.TotalSecondsActive,
		HoursActive:        int(p.TotalSecondsActive / 3600),
	}
}

type AchievementResponse struct {
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Tier        int16      `json:"tier"`
	Metric      string     `json:"metric"`
	Threshold   int        `json:"threshold"`
	XPReward    int        `json:"xp_reward"`
	Unlocked    bool       `json:"unlocked"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
}

func toAchievementResponse(a AchievementStatus) AchievementResponse {
	return AchievementResponse{
		Code:        a.Code,
		Title:       a.Title,
		Description: a.Description,
		Tier:        a.Tier,
		Metric:      a.Metric,
		Threshold:   a.Threshold,
		XPReward:    a.XPReward,
		Unlocked:    a.Unlocked,
		UnlockedAt:  a.UnlockedAt,
	}
}
