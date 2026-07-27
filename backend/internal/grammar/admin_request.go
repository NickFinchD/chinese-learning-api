package grammar

type AdminListRequest struct {
	Search string `form:"search"`
	HSK    int16  `form:"hsk"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type CreateNoteRequest struct {
	Title               string `json:"title" binding:"required"`
	Explanation         string `json:"explanation" binding:"required"`
	ExampleHanzi        string `json:"example_hanzi"`
	ExamplePinyin       string `json:"example_pinyin"`
	ExampleTranslation  string `json:"example_translation"`
	Example2Hanzi       string `json:"example2_hanzi"`
	Example2Pinyin      string `json:"example2_pinyin"`
	Example2Translation string `json:"example2_translation"`
	Example3Hanzi       string `json:"example3_hanzi"`
	Example3Pinyin      string `json:"example3_pinyin"`
	Example3Translation string `json:"example3_translation"`
	HSKLevel            int16  `json:"hsk_level" binding:"required"`
}

type UpdateNoteRequest struct {
	Title               string `json:"title" binding:"required"`
	Explanation         string `json:"explanation" binding:"required"`
	ExampleHanzi        string `json:"example_hanzi"`
	ExamplePinyin       string `json:"example_pinyin"`
	ExampleTranslation  string `json:"example_translation"`
	Example2Hanzi       string `json:"example2_hanzi"`
	Example2Pinyin      string `json:"example2_pinyin"`
	Example2Translation string `json:"example2_translation"`
	Example3Hanzi       string `json:"example3_hanzi"`
	Example3Pinyin      string `json:"example3_pinyin"`
	Example3Translation string `json:"example3_translation"`
	HSKLevel            int16  `json:"hsk_level" binding:"required"`
}
