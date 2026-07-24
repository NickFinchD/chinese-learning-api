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
