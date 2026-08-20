package middleware

import (
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	tokenService *TokenService
}

type contextKey string

const UserIDKey contextKey = "userID"

func NewAuthMiddleware(tokenService *TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Отсутствует Authorization", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Неверный формат Authorization", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Недействительный token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
