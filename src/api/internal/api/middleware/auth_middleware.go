package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type AuthMiddleware struct {
	identityManager domain.IdentityManager
}

func NewAuthMiddleware(im domain.IdentityManager) *AuthMiddleware {
	return &AuthMiddleware{identityManager: im}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
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
