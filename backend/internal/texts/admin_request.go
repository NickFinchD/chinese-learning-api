package texts

type AdminListRequest struct {
	Search string `form:"search"`
	HSK    int16  `form:"hsk"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type CreateTextRequest struct {
	Title       string `json:"title" binding:"required"`
	Hanzi       string `json:"hanzi" binding:"required"`
	Pinyin      string `json:"pinyin" binding:"required"`
	Translation string `json:"translation" binding:"required"`
	HSKLevel    int16  `json:"hsk_level" binding:"required"`
}

type UpdateTextRequest struct {
	Title       string `json:"title" binding:"required"`
	Hanzi       string `json:"hanzi" binding:"required"`
	Pinyin      string `json:"pinyin" binding:"required"`
	Translation string `json:"translation" binding:"required"`
	HSKLevel    int16  `json:"hsk_level" binding:"required"`
}
