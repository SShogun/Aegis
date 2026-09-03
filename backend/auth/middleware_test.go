package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sshogun/Aegis/backend/users"
)

type fakeUserService struct {
	user *users.User
	err  error
}

func (s *fakeUserService) FindOrCreateFromClaims(
	ctx context.Context,
	provider string,
	subject string,
	email string,
	name string,
) (*users.User, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeUserService) FindByID(ctx context.Context, userID string) (*users.User, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.user, nil
}

type fakeSessionService struct {
	session          *users.Session
	err              error
	deleteErr        error
	deletedSessionID string
	deleteCalls      int
}

func (s *fakeSessionService) Create(ctx context.Context, userID string) (*users.Session, error) {
	return nil, errors.New("not implemented")
}

func (s *fakeSessionService) FindValid(
	ctx context.Context,
	sessionID string,
) (*users.Session, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.session, nil
}

func (s *fakeSessionService) Delete(ctx context.Context, sessionID string) error {
	s.deleteCalls++
	s.deletedSessionID = sessionID
	return s.deleteErr
}

func TestRequireAuthenticationReturnsUnauthorizedWithoutSessionCookie(t *testing.T) {
	handler := NewHandler(nil, &fakeUserService{}, &fakeSessionService{})

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()

	handler.RequireAuthentication(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}
}

func TestRequireAuthenticationReturnsUnauthorizedForInvalidSession(t *testing.T) {
	handler := NewHandler(
		nil,
		&fakeUserService{},
		&fakeSessionService{err: users.ErrSessionNotFound},
	)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "missing-session"})
	rr := httptest.NewRecorder()

	handler.RequireAuthentication(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}
}

func TestRequireAuthenticationAddsPrincipalToContext(t *testing.T) {
	expectedUser := &users.User{
		ID:               "user-1",
		Email:            "user@example.com",
		Name:             "Test User",
		IdentityProvider: "test-provider",
	}

	handler := NewHandler(
		nil,
		&fakeUserService{user: expectedUser},
		&fakeSessionService{
			session: &users.Session{
				ID:        "session-1",
				UserID:    expectedUser.ID,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
	)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Fatal("expected principal in request context")
		}

		if principal.UserID != expectedUser.ID {
			t.Fatalf("expected principal user ID %q, got %q", expectedUser.ID, principal.UserID)
		}

		if principal.Email != expectedUser.Email {
			t.Fatalf("expected principal email %q, got %q", expectedUser.Email, principal.Email)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-1"})
	rr := httptest.NewRecorder()

	handler.RequireAuthentication(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestMeReturnsAuthenticatedIdentity(t *testing.T) {
	handler := NewHandler(nil, &fakeUserService{}, &fakeSessionService{})

	principal := Principal{
		UserID:           "user-1",
		Email:            "user@example.com",
		Name:             "Test User",
		IdentityProvider: "test-provider",
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = req.WithContext(ContextWithPrincipal(req.Context(), principal))
	rr := httptest.NewRecorder()

	handler.Me(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := map[string]string{
		"id":                principal.UserID,
		"email":             principal.Email,
		"name":              principal.Name,
		"identity_provider": principal.IdentityProvider,
	}

	for key, value := range expected {
		if response[key] != value {
			t.Fatalf("expected response[%q] %q, got %q", key, value, response[key])
		}
	}

	if _, ok := response["provider_subject"]; ok {
		t.Fatal("expected provider_subject not to be exposed")
	}

	if _, ok := response["session_id"]; ok {
		t.Fatal("expected session_id not to be exposed")
	}
}

func TestLogoutDeletesSessionAndClearsCookie(t *testing.T) {
	sessionService := &fakeSessionService{}
	handler := NewHandler(nil, &fakeUserService{}, sessionService)

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-1"})
	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if sessionService.deleteCalls != 1 {
		t.Fatalf("expected one session deletion, got %d", sessionService.deleteCalls)
	}
	if sessionService.deletedSessionID != "session-1" {
		t.Fatalf("expected session-1 to be deleted, got %q", sessionService.deletedSessionID)
	}
	assertLogoutResponse(t, rr, http.StatusOK, "logged out")
}

func TestLogoutWithoutSessionCookieStillClearsCookie(t *testing.T) {
	sessionService := &fakeSessionService{}
	handler := NewHandler(nil, &fakeUserService{}, sessionService)

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if sessionService.deleteCalls != 0 {
		t.Fatalf("expected no session deletion, got %d calls", sessionService.deleteCalls)
	}
	assertLogoutResponse(t, rr, http.StatusOK, "logged out")
}

func TestLogoutReturnsGenericErrorAndClearsCookieWhenDeletionFails(t *testing.T) {
	sessionService := &fakeSessionService{deleteErr: errors.New("database unavailable")}
	handler := NewHandler(nil, &fakeUserService{}, sessionService)

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-1"})
	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if sessionService.deleteCalls != 1 {
		t.Fatalf("expected one session deletion, got %d", sessionService.deleteCalls)
	}
	assertLogoutResponse(t, rr, http.StatusInternalServerError, "logout failed")
	if strings.Contains(rr.Body.String(), "database unavailable") {
		t.Fatal("expected database error not to be exposed")
	}
}

func assertLogoutResponse(t *testing.T, rr *httptest.ResponseRecorder, status int, message string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode logout response: %v", err)
	}
	if response.Message != message {
		t.Fatalf("expected logout message %q, got %q", message, response.Message)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one clearing cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != "session_id" || cookie.Value != "" || cookie.Path != "/" ||
		!cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode || cookie.MaxAge != -1 {
		t.Fatalf("unexpected clearing cookie: %#v", cookie)
	}
}
