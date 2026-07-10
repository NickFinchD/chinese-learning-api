package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	existingUser, err := s.repository.GetByEmail(ctx, req.Email)

	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	err = s.repository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
