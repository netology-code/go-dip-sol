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

type CommentHandler struct {
	commentService service.CommentServiceInterface
	eventLogger    *logger.EventLogger
}

func NewCommentHandler(commentService service.CommentServiceInterface, eventLogger *logger.EventLogger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		eventLogger:    eventLogger,
	}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		WriteError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		WriteError(w, "Content is required", http.StatusBadRequest)
		return
	}
	if len(req.Content) > 1000 {
		WriteError(w, "Content exceeds maximum length of 1000 characters", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.Create(r.Context(), userID, postID, req.Content)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	h.eventLogger.LogEvent(fmt.Sprintf("user %d created comment %d", userID, comment.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		WriteError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	limit := 20
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

	comments, total, err := h.commentService.GetByPost(r.Context(), postID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	resp := struct {
		Comments []*model.Comment `json:"comments"`
		Total    int              `json:"total"`
		Limit    int              `json:"limit"`
		Offset   int              `json:"offset"`
		PostID   int              `json:"post_id"`
	}{
		Comments: comments,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
		PostID:   postID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
