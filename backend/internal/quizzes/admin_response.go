package quizzes

import "time"

// AdminOptionResponse mirrors Option but exposes is_correct and sort_order,
// which the public Option hides (json:"-") so players can't read the
// answer out of the API response. Admin routes are behind RequireAdmin, so
// exposing them here is safe and necessary for editing.
type AdminOptionResponse struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Pinyin    string `json:"pinyin,omitempty"`
	IsCorrect bool   `json:"is_correct"`
	SortOrder int16  `json:"sort_order"`
}

type AdminQuizResponse struct {
	ID        int64                 `json:"id"`
	Question  string                `json:"question"`
	HSKLevel  int16                 `json:"hsk_level"`
	Direction string                `json:"direction"`
	Hanzi     string                `json:"hanzi,omitempty"`
	Pinyin    string                `json:"pinyin,omitempty"`
	Options   []AdminOptionResponse `json:"options"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

func toAdminQuizResponse(q Quiz) AdminQuizResponse {

	options := make([]AdminOptionResponse, 0, len(q.Options))

	for _, o := range q.Options {

		pinyin := ""
		if o.Pinyin != nil {
			pinyin = *o.Pinyin
		}

		options = append(options, AdminOptionResponse{
			ID:        o.ID,
			Text:      o.Text,
			Pinyin:    pinyin,
			IsCorrect: o.IsCorrect,
			SortOrder: o.SortOrder,
		})
	}

	hanzi := ""
	if q.Hanzi != nil {
		hanzi = *q.Hanzi
	}

	pinyin := ""
	if q.Pinyin != nil {
		pinyin = *q.Pinyin
	}

	return AdminQuizResponse{
		ID:        q.ID,
		Question:  q.Question,
		HSKLevel:  q.HSKLevel,
		Direction: q.Direction,
		Hanzi:     hanzi,
		Pinyin:    pinyin,
		Options:   options,
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
	}
}

func toAdminQuizResponses(quizzes []Quiz) []AdminQuizResponse {

	result := make([]AdminQuizResponse, 0, len(quizzes))

	for _, q := range quizzes {
		result = append(result, toAdminQuizResponse(q))
	}

	return result
}
