package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type AuthMiddleware struct {
	identityManager domain.IdentityManager
}

func NewAuthMiddleware(im domain.IdentityManager) *AuthMiddleware {
	return &AuthMiddleware{identityManager: im}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		if token == "" {
			cookie, err := r.Cookie("nebula_token")
			if err == nil && cookie.Value != "" {
				token = cookie.Value
			}
		}

		if token == "" {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}

		user, err := m.identityManager.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Inject user into context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) AuthenticateRequest(r *http.Request) (*domain.User, error) {
	token := ""
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		}
	}
	if token == "" {
		cookie, err := r.Cookie("nebula_token")
		if err == nil && cookie.Value != "" {
			token = cookie.Value
		}
	}
	if token == "" {
		return nil, http.ErrNoCookie
	}
	return m.identityManager.ValidateToken(r.Context(), token)
}
