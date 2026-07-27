package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestContext(userID int64, setUserID bool) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	if setUserID {
		c.Set("user_id", userID)
	}

	return c, recorder
}

func TestRequireAdmin_AllowsAdmin(t *testing.T) {
	repo := &fakeRepository{isAdminResult: true}
	service := NewService(repo, testConfig())

	c, _ := newTestContext(1, true)

	RequireAdmin(service)(c)

	if c.IsAborted() {
		t.Fatal("expected the request not to be aborted for an admin user")
	}
}

func TestRequireAdmin_BlocksNonAdmin(t *testing.T) {
	repo := &fakeRepository{isAdminResult: false}
	service := NewService(repo, testConfig())

	c, recorder := newTestContext(1, true)

	RequireAdmin(service)(c)

	if !c.IsAborted() {
		t.Fatal("expected the request to be aborted for a non-admin user")
	}

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}
}

func TestRequireAdmin_BlocksOnLookupError(t *testing.T) {
	repo := &fakeRepository{isAdminErr: errors.New("connection refused")}
	service := NewService(repo, testConfig())

	c, recorder := newTestContext(1, true)

	RequireAdmin(service)(c)

	if !c.IsAborted() {
		t.Fatal("expected the request to be aborted when the admin lookup fails")
	}

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
}

func TestRequireAdmin_BlocksMissingUserID(t *testing.T) {
	// GetUserID falls back to 0 when user_id was never set. A real repository
	// finds no matching row for id 0 and reports IsAdmin=false (see
	// Repository.IsAdmin's ErrNoRows handling) — mirror that here.
	repo := &fakeRepository{isAdminResult: false}
	service := NewService(repo, testConfig())

	c, recorder := newTestContext(0, false)

	RequireAdmin(service)(c)

	if !c.IsAborted() {
		t.Fatal("expected the request to be aborted when user_id is missing")
	}

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for user_id 0 (no such user), got %d", recorder.Code)
	}
}
