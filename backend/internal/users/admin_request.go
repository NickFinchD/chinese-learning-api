package users

type AdminListRequest struct {
	Search string `form:"search"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type SetAdminRequest struct {
	IsAdmin bool `json:"is_admin"`
}
