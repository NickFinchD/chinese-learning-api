package learning

import "time"

type ProgressResponse struct {
	WordID         int64      `json:"word_id"`
	Stage          int16      `json:"stage"`
	MaxStage       int16      `json:"max_stage"`
	NextEligibleAt *time.Time `json:"next_eligible_at,omitempty"`
	Learned        bool       `json:"learned"`
}

func toProgressResponse(progress *WordLearningProgress) ProgressResponse {
	return ProgressResponse{
		WordID:         progress.WordID,
		Stage:          progress.Stage,
		MaxStage:       MaxStage,
		NextEligibleAt: progress.NextEligibleAt,
		Learned:        progress.LearnedAt != nil,
	}
}

type InProgressWordResponse struct {
	WordID          int64      `json:"word_id"`
	Hanzi           string     `json:"hanzi"`
	Pinyin          string     `json:"pinyin"`
	Translation     string     `json:"translation"`
	HSKLevel        int16      `json:"hsk_level"`
	Stage           int16      `json:"stage"`
	MaxStage        int16      `json:"max_stage"`
	RepetitionsLeft int16      `json:"repetitions_left"`
	NextEligibleAt  *time.Time `json:"next_eligible_at,omitempty"`
}

func toInProgressResponse(detail WordProgressDetail) InProgressWordResponse {
	return InProgressWordResponse{
		WordID:          detail.WordID,
		Hanzi:           detail.Hanzi,
		Pinyin:          detail.Pinyin,
		Translation:     detail.Translation,
		HSKLevel:        detail.HSKLevel,
		Stage:           detail.Stage,
		MaxStage:        MaxStage,
		RepetitionsLeft: MaxStage - detail.Stage,
		NextEligibleAt:  detail.NextEligibleAt,
	}
}
