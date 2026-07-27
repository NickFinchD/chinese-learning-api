package sentences

type AdminListRequest struct {
	Search string `form:"search"`
	HSK    int16  `form:"hsk"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type AdminExerciseRequest struct {
	Translation string   `json:"translation" binding:"required"`
	Chunks      []string `json:"chunks" binding:"required,min=2"`
	Pinyin      string   `json:"pinyin"`
	HSKLevel    int16    `json:"hsk_level" binding:"required"`
}
