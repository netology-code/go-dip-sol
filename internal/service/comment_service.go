package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"context"
	"fmt"
	"strings"
)

type CommentService struct {
	repo       repository.CommentRepository
	postRepo   repository.PostRepository
}

func NewCommentService(repo repository.CommentRepository, postRepo repository.PostRepository) *CommentService {
	return &CommentService{
		repo:     repo,
		postRepo: postRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, userID, postID int, content string) (*model.Comment, error) {
	if postID <= 0 {
		return nil, apperrors.ErrInvalidPostID
	}

	exists, err := s.postRepo.Exists(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !exists {
		return nil, apperrors.ErrPostNotFound
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}
	if len(content) > 1000 {
		return nil, fmt.Errorf("content exceeds 1000 characters")
	}

	comment := &model.Comment{
		PostID:   postID,
		AuthorID: userID,
		Content:  content,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

func (s *CommentService) GetByPost(ctx context.Context, postID, limit, offset int) ([]*model.Comment, int, error) {
	if postID <= 0 {
		return nil, 0, apperrors.ErrInvalidPostID
	}

	exists, err := s.postRepo.Exists(ctx, postID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !exists {
		return nil, 0, apperrors.ErrPostNotFound
	}

	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	comments, err := s.repo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}

	total, err := s.repo.GetCountByPostID(ctx, postID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return comments, total, nil
}
