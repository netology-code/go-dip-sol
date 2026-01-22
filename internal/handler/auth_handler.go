package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// AuthHandler обрабатывает HTTP-запросы для аутентификации пользователей
type AuthHandler struct {
	userService service.UserServiceInterface // Интерфейс - для поддержки моков в тестах
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(userService service.UserServiceInterface) *AuthHandler { // принимает интерфейс
	return &AuthHandler{
		userService: userService,
	}
}

// Register обрабатывает регистрацию нового пользователя
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserCreateRequest true "User registration data"
// @Success 201 {object} model.TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenResp, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			h.respondWithError(w, "Validation error: invalid input", http.StatusBadRequest)
			return
		}

		switch err {
		case apperrors.ErrUserAlreadyExists:
			h.respondWithError(w, "User already exists", http.StatusConflict)
		default:
			h.respondWithError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.respondWithJSON(w, tokenResp, http.StatusCreated)
}

// Login обрабатывает вход пользователя в систему
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.UserLoginRequest true "User login data"
// @Success 200 {object} model.TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenResp, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			h.respondWithError(w, "Validation error: invalid input", http.StatusBadRequest)
			return
		}

		switch err {
		case apperrors.ErrInvalidCredentials:
			h.respondWithError(w, "Invalid email or password", http.StatusUnauthorized)
		default:
			h.respondWithError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.respondWithJSON(w, tokenResp, http.StatusOK)
}

// GetProfile получает профиль текущего пользователя
// @Summary Get user profile
// @Description Get the profile of the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case apperrors.ErrUserNotFound:
			h.respondWithError(w, "User not found", http.StatusNotFound)
		default:
			h.respondWithError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Преобразование в UserResponse для ответа
	userResp := model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	h.respondWithJSON(w, userResp, http.StatusOK)
}

// respondWithError отправляет JSON-ответ с ошибкой
func (h *AuthHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON отправляет JSON-ответ с данными
func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
