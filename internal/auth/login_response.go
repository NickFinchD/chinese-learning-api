package auth

type LoginResponse struct {
	Token string           `json:"token"`
	User  RegisterResponse `json:"user"`
}
