package middleware

import (
	"advanced-blog-management-system/pkg/auth"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// contextKey - кастомный тип, чтобы избежать коллизии
type contextKey string

const (
	// UserIDKey - для сохранения user ID в контексте
	UserIDKey contextKey = "userID"
	// UserEmailKey - ключ для сохранения email пользователя в контексте
	UserEmailKey contextKey = "userEmail"
	// UserNameKey - ключ для сохранения username в контекс
	UserNameKey contextKey = "username"
)

// AuthMiddleware обеспечивает JWT аутентификацию
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// NewAuthMiddleware создает новый инстанс auth middleware
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth - middleware требует валидный JWT token
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Извлечь токен из заголовка Authorization (Bearer токен)
		token := extractToken(r)
		if token == "" {
			writeJSONError(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// 2. Валидировать токен через jwtManager
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			writeJSONError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 3. Добавить данные пользователя в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserNameKey, claims.Username)

		// 4. Передать управление следующему handler
		next(w, r.WithContext(ctx))
	}
}

// OptionalAuth - middleware извлекает JWT token, если есть, но не требует его
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Попытаться извлечь токен из заголовка
		token := extractToken(r)
		if token == "" {
			next(w, r)
			return
		}

		// 2. Если токен есть, то валидировать его
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			// 4. Если токен невалидный, то продолжить как анонимный
			next(w, r)
			return
		}

		// 3. Если токен валидный, то добавить данные в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserNameKey, claims.Username)

		// 5. Передать управление следующему handler
		next(w, r.WithContext(ctx))
	}
}

// extractToken извлекает JWT токен из заголовка Authorization
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Проверяем формат "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// GetUserIDFromContext - извлекает ID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// GetUserEmailFromContext извлекает email пользователя из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// GetUsernameFromContext извлекает username из контекста
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UserNameKey).(string)
	return username, ok
}

// writeJSONError отправляет ошибку в формате JSON
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Chain позволяет объединить несколько middleware в цепочку
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
