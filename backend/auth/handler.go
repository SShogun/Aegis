package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sshogun/Aegis/backend/users"
)

type userService interface {
	FindOrCreateFromClaims(
		ctx context.Context,
		provider string,
		subject string,
		email string,
		name string,
	) (*users.User, error)
	FindByID(ctx context.Context, userID string) (*users.User, error)
}

type sessionService interface {
	Create(ctx context.Context, userID string) (*users.Session, error)
	FindValid(ctx context.Context, sessionID string) (*users.Session, error)
	Delete(ctx context.Context, sessionID string) error
}

type Handler struct {
	OIDC           *OIDC
	UserService    userService
	SessionService sessionService
}

func NewHandler(
	oidcClient *OIDC,
	userService userService,
	sessionService sessionService,
) *Handler {
	return &Handler{
		OIDC:           oidcClient,
		UserService:    userService,
		SessionService: sessionService,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := GenerateState()
	if err != nil {
		http.Error(
			w,
			"failed to generate authentication state",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})

	authorizationURL := h.OIDC.AuthorizationURL(state)

	http.Redirect(w, r, authorizationURL, http.StatusFound)
}

func (h *Handler) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil || sessionCookie.Value == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := h.SessionService.FindValid(r.Context(), sessionCookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := h.UserService.FindByID(r.Context(), session.UserID)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		principal := Principal{
			UserID:           user.ID,
			Email:            user.Email,
			Name:             user.Name,
			IdentityProvider: user.IdentityProvider,
		}

		next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), principal)))
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := struct {
		ID               string `json:"id"`
		Email            string `json:"email"`
		Name             string `json:"name"`
		IdentityProvider string `json:"identity_provider"`
	}{
		ID:               principal.UserID,
		Email:            principal.Email,
		Name:             principal.Name,
		IdentityProvider: principal.IdentityProvider,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		if err := h.SessionService.Delete(r.Context(), sessionCookie.Value); err != nil {
			http.SetCookie(w, expiredSessionCookie())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
			}{Message: "logout failed"})
			return
		}
	}

	http.SetCookie(w, expiredSessionCookie())
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: "logged out"})
}

func expiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state == "" || code == "" {
		http.Error(
			w,
			"missing authentication callback parameters",
			http.StatusBadRequest,
		)
		return
	}

	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(
			w,
			"authentication state not found",
			http.StatusBadRequest,
		)
		return
	}

	if state != stateCookie.Value {
		http.Error(
			w,
			"invalid authentication state",
			http.StatusBadRequest,
		)
		return
	}

	claims, err := h.OIDC.ExchangeCode(r.Context(), code)
	if err != nil {
		http.Error(
			w,
			"authentication failed",
			http.StatusUnauthorized,
		)
		return
	}

	user, err := h.UserService.FindOrCreateFromClaims(
		r.Context(),
		h.OIDC.ProviderName,
		claims.Subject,
		claims.Email,
		claims.Name,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create or find user",
			http.StatusInternalServerError,
		)
		return
	}

	// Create a persistent login session for the authenticated user.
	session, err := h.SessionService.Create(
		r.Context(),
		user.ID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create session",
			http.StatusInternalServerError,
		)
		return
	}

	// The OAuth state is no longer needed.
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// Give the browser the session ID so it can identify
	// this authenticated session on future requests.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(users.SessionDuration.Seconds()),
	})

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(
		"authenticated as " + user.Email,
	))
}
