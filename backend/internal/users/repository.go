package users

import (
	"context"
	"strconv"

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

func (r *Repository) AdminList(ctx context.Context, request AdminListRequest) ([]User, int, error) {

	page := request.Page
	if page < 1 {
		page = 1
	}

	limit := request.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	where := `WHERE 1=1`
	args := []interface{}{}
	arg := 1

	if request.Search != "" {
		where += `
			AND (
				username ILIKE $` + strconv.Itoa(arg) + `
				OR email ILIKE $` + strconv.Itoa(arg) + `
			)
		`

		args = append(args, "%"+request.Search+"%")
		arg++
	}

	var total int

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, username, email, is_admin, created_at
		FROM users
	` + where + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	result := make([]User, 0)

	for rows.Next() {

		var u User

		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.IsAdmin, &u.CreatedAt); err != nil {
			return nil, 0, err
		}

		result = append(result, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {

	u := &User{}

	err := r.db.QueryRow(ctx, `
		SELECT id, username, email, is_admin, created_at FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Username, &u.Email, &u.IsAdmin, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) CountAdmins(ctx context.Context) (int, error) {

	var count int

	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE is_admin = true`).Scan(&count)

	return count, err
}

func (r *Repository) SetAdmin(ctx context.Context, id int64, isAdmin bool) (*User, error) {

	u := &User{}

	err := r.db.QueryRow(ctx, `
		UPDATE users SET is_admin = $1 WHERE id = $2
		RETURNING id, username, email, is_admin, created_at
	`, isAdmin, id).Scan(&u.ID, &u.Username, &u.Email, &u.IsAdmin, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil
}
