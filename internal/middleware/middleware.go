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

const RoleKey contextKey = "role"

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

		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Недействительный token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value(RoleKey).(string)
		if !ok {
			http.Error(w, "у пользователя нет роли", http.StatusUnauthorized)
			return
		}

		if userRole != role {
			http.Error(w, "недостаточно прав", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
