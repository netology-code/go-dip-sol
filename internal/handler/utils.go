package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse - структура для ответа с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// WriteError отправляет JSON-ответ с ошибкой
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

// HandleServiceError обрабатывает ошибки сервиса и отправляет соответствующий ответ
func HandleServiceError(w http.ResponseWriter, err error) {
	if _, ok := err.(validator.ValidationErrors); ok {
		WriteError(w, "Validation error: invalid input", http.StatusBadRequest)
		return
	}

	switch {
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		WriteError(w, "User already exists", http.StatusConflict)
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		WriteError(w, "Invalid email or password", http.StatusUnauthorized)
	case errors.Is(err, apperrors.ErrPostNotFound):
		WriteError(w, "Post not found", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrCommentNotFound):
		WriteError(w, "Comment not found", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrForbidden):
		WriteError(w, "Forbidden", http.StatusForbidden)
	case errors.Is(err, apperrors.ErrUnauthorized):
		WriteError(w, "Unauthorized", http.StatusUnauthorized)
	default:
		WriteError(w, "Internal server error", http.StatusInternalServerError)
	}
}
