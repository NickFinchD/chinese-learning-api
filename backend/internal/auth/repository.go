package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}
func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			username,
			email,
			password_hash
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)
}
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password_hash,
			avatar,
			is_admin,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Avatar,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
// GetByID returns (nil, nil) — not an error — when no such user exists, so
// callers (e.g. /me for a JWT whose user was since deleted) can tell "no
// longer a valid session" apart from a real database failure.
func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password_hash,
			avatar,
			is_admin,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Avatar,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// IsAdmin is a cheap single-column check used by RequireAdmin on every
// admin request, so demoting an admin takes effect on their very next
// request instead of waiting for their JWT to expire (see middleware.go).
func (r *Repository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT is_admin FROM users WHERE id = $1`

	var isAdmin bool

	err := r.db.QueryRow(ctx, query, userID).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return isAdmin, nil
}
