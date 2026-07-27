package words

type AdminListRequest struct {
	Search string `form:"search"`
	HSK    int16  `form:"hsk"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type ExampleRequest struct {
	Hanzi       string `json:"hanzi" binding:"required"`
	Pinyin      string `json:"pinyin" binding:"required"`
	Translation string `json:"translation" binding:"required"`
}

type CreateWordRequest struct {
	Hanzi        string           `json:"hanzi" binding:"required"`
	Pinyin       string           `json:"pinyin" binding:"required"`
	Translation  string           `json:"translation" binding:"required"`
	PartOfSpeech string           `json:"part_of_speech" binding:"required"`
	HSKLevel     int16            `json:"hsk_level" binding:"required"`
	Examples     []ExampleRequest `json:"examples"`
}

type UpdateWordRequest struct {
	Hanzi        string           `json:"hanzi" binding:"required"`
	Pinyin       string           `json:"pinyin" binding:"required"`
	Translation  string           `json:"translation" binding:"required"`
	PartOfSpeech string           `json:"part_of_speech" binding:"required"`
	HSKLevel     int16            `json:"hsk_level" binding:"required"`
	Examples     []ExampleRequest `json:"examples"`
}
