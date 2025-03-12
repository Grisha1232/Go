package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-app/internal/services"
)

// AuthMiddleware проверяет JWT токен и передает userID в контекст запроса
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := services.ValidateToken(token)
		if err != nil {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Передаем userID в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}
