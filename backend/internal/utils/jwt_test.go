package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {

	token, err := GenerateToken(42, "test-secret")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := ParseToken(token, "test-secret")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
	}
}

func TestParseToken_WrongSecret(t *testing.T) {

	token, err := GenerateToken(1, "correct-secret")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ParseToken(token, "wrong-secret")

	if err == nil {
		t.Fatal("expected an error for a token signed with a different secret")
	}
}

func TestParseToken_Expired(t *testing.T) {

	claims := Claims{
		UserID: 7,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := raw.SignedString([]byte("test-secret"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ParseToken(token, "test-secret")

	if err == nil {
		t.Fatal("expected an error for an expired token")
	}
}
