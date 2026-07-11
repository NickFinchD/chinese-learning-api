package auth

type LoginResult struct {
	User  *User
	Token string
}
