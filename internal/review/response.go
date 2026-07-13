package review

type ReviewWordResponse struct {
	WordID      int64  `json:"word_id"`
	Hanzi       string `json:"hanzi"`
	Pinyin      string `json:"pinyin"`
	Translation string `json:"translation"`
}
