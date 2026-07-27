package users

import "time"

// User is the admin-facing view of an account — deliberately excludes
// PasswordHash (that lives only in auth.User, never serialized here).
type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
