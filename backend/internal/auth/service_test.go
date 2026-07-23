package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	getByEmailUser *User
	getByEmailErr  error

	getByIDUser *User
	getByIDErr  error

	createErr     error
	createCalled  bool
	createdUser   *User
}

func (f *fakeRepository) Create(ctx context.Context, user *User) error {
	f.createCalled = true
	f.createdUser = user

	if f.createErr != nil {
		return f.createErr
	}

	user.ID = 1

	return nil
}

func (f *fakeRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return f.getByEmailUser, f.getByEmailErr
}

func (f *fakeRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	return f.getByIDUser, f.getByIDErr
}

func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}
}

func TestRegister_Success(t *testing.T) {

	repo := &fakeRepository{
		getByEmailErr: pgx.ErrNoRows,
	}

	service := NewService(repo, testConfig())

	user, err := service.Register(context.Background(), RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "supersecret",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repository.Create to be called")
	}

	if user.PasswordHash == "supersecret" {
		t.Fatal("password must be hashed before storing")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("supersecret")) != nil {
		t.Fatal("stored hash does not match the original password")
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {

	repo := &fakeRepository{
		getByEmailUser: &User{ID: 1, Email: "alice@example.com"},
	}

	service := NewService(repo, testConfig())

	_, err := service.Register(context.Background(), RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "supersecret",
	})

	if err == nil {
		t.Fatal("expected an error for a duplicate email")
	}

	if repo.createCalled {
		t.Fatal("repository.Create should not be called when the email is already taken")
	}
}

func TestRegister_PropagatesUnexpectedLookupError(t *testing.T) {

	dbErr := errors.New("connection refused")

	repo := &fakeRepository{
		getByEmailErr: dbErr,
	}

	service := NewService(repo, testConfig())

	_, err := service.Register(context.Background(), RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "supersecret",
	})

	if !errors.Is(err, dbErr) {
		t.Fatalf("expected the lookup error to propagate, got %v", err)
	}

	if repo.createCalled {
		t.Fatal("repository.Create should not be called when the lookup fails unexpectedly")
	}
}

func TestLogin_Success(t *testing.T) {

	hash, err := bcrypt.GenerateFromPassword([]byte("supersecret"), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo := &fakeRepository{
		getByEmailUser: &User{ID: 1, Email: "alice@example.com", PasswordHash: string(hash)},
	}

	service := NewService(repo, testConfig())

	result, err := service.Login(context.Background(), LoginRequest{
		Email:    "alice@example.com",
		Password: "supersecret",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Token == "" {
		t.Fatal("expected a non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {

	hash, err := bcrypt.GenerateFromPassword([]byte("supersecret"), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo := &fakeRepository{
		getByEmailUser: &User{ID: 1, Email: "alice@example.com", PasswordHash: string(hash)},
	}

	service := NewService(repo, testConfig())

	_, err = service.Login(context.Background(), LoginRequest{
		Email:    "alice@example.com",
		Password: "wrong-password",
	})

	if err == nil {
		t.Fatal("expected an error for a wrong password")
	}
}

func TestLogin_UserNotFound(t *testing.T) {

	repo := &fakeRepository{
		getByEmailErr: pgx.ErrNoRows,
	}

	service := NewService(repo, testConfig())

	_, err := service.Login(context.Background(), LoginRequest{
		Email:    "nobody@example.com",
		Password: "whatever",
	})

	if err == nil {
		t.Fatal("expected an error when the user does not exist")
	}
}

func TestMe_ReturnsRepositoryResult(t *testing.T) {

	repo := &fakeRepository{
		getByIDUser: &User{ID: 5, Username: "bob"},
	}

	service := NewService(repo, testConfig())

	user, err := service.Me(context.Background(), 5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Username != "bob" {
		t.Fatalf("expected username bob, got %s", user.Username)
	}
}
