package auth

import (
	"context"
	"errors"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/utils"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
	config     *config.Config
}

func NewService(repository *Repository, cfg *config.Config) *Service {
	return &Service{
		repository: repository,
		config:     cfg,
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
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResult, error) {

	user, err := s.repository.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, s.config.JWT.Secret)

	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}
func (s *Service) Me(ctx context.Context, userID int64) (*User, error) {

	user, err := s.repository.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}
