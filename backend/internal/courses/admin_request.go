package courses

type AdminListRequest struct {
	Search string `form:"search"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type CourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	HSKLevel    int16  `json:"hsk_level" binding:"required"`
	SortOrder   int    `json:"sort_order"`
}
