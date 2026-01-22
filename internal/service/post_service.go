package service

import (
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"context"
	"fmt"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *PostService) GetByID(ctx context.Context, id int, requestorID int) (*model.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	total, err := s.postRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total post count: %w", err)
	}

	return posts, total, nil
}

func (s *PostService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.GetByAuthorID(ctx, authorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts by author: %w", err)
	}

	total, err := s.postRepo.GetTotalCountByAuthorID(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total post count by author: %w", err)
	}

	return posts, total, nil
}
