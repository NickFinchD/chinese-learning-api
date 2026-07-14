package review

type ReviewWordResponse struct {
	WordID      int64  `json:"word_id"`
	Hanzi       string `json:"hanzi"`
	Pinyin      string `json:"pinyin"`
	Translation string `json:"translation"`
}
type StatisticsResponse struct {
	TotalWords     int `json:"total_words"`
	ReadyForReview int `json:"ready_for_review"`
	ReviewedWords  int `json:"reviewed_words"`
}
type ReviewSessionResponse struct {
	Total int                  `json:"total"`
	Words []ReviewWordResponse `json:"words"`
}
