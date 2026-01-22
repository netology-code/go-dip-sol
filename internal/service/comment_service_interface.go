package service

import (
	"advanced-blog-management-system/internal/model"
	"context"
)

type CommentServiceInterface interface {
	Create(ctx context.Context, userID, postID int, content string) (*model.Comment, error)

	GetByPost(ctx context.Context, postID, limit, offset int) ([]*model.Comment, int, error)
}
