package service

import (
	"advanced-blog-management-system/internal/model"
	"context"
)

type PostServiceInterface interface {
	Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error)

	GetByID(ctx context.Context, id int, requestorID int) (*model.Post, error)

	GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error)

	GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error)
}
