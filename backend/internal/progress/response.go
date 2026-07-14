package progress

type ProgressResponse struct {
	Status      string `json:"status"`
	CurrentStep int    `json:"current_step"`
	Score       int    `json:"score"`
}
