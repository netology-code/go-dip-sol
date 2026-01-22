package handler

import (
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService service.PostServiceInterface
	eventLogger *logger.EventLogger
}

func NewPostHandler(postService service.PostServiceInterface, eventLogger *logger.EventLogger) *PostHandler {
	return &PostHandler{
		postService: postService,
		eventLogger: eventLogger,
	}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postService.Create(r.Context(), userID, &req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	h.eventLogger.LogEvent(fmt.Sprintf("user %d created post %d", userID, post.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var requestorID int
	if userID, ok := middleware.GetUserIDFromContext(r.Context()); ok {
		requestorID = userID
	}

	post, err := h.postService.GetByID(r.Context(), id, requestorID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, total, err := h.postService.GetAll(r.Context(), limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	type PostsResponse struct {
		Posts  []*model.Post `json:"posts"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	resp := PostsResponse{
		Posts:  posts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authorIDStr := chi.URLParam(r, "authorID")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		WriteError(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, total, err := h.postService.GetByAuthor(r.Context(), authorID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	type PostsResponse struct {
		Posts    []*model.Post `json:"posts"`
		Total    int           `json:"total"`
		Limit    int           `json:"limit"`
		Offset   int           `json:"offset"`
		AuthorID int           `json:"author_id"`
	}

	resp := PostsResponse{
		Posts:    posts,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
		AuthorID: authorID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
