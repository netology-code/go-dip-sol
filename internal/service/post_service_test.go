package service

import (
	"advanced-blog-management-system/internal/model"

	"context"
	"errors"
	"testing"
	"time"
)

// mockPostRepo is a mock implementation of PostRepository
type mockPostRepo struct {
	createFunc                   func(ctx context.Context, post *model.Post) error
	getByIDFunc                  func(ctx context.Context, id int) (*model.Post, error)
	getAllFunc                   func(ctx context.Context, limit, offset int) ([]*model.Post, error)
	getTotalCountFunc            func(ctx context.Context) (int, error)
	existsFunc                   func(ctx context.Context, id int) (bool, error)
	getByAuthorIDFunc            func(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error)
	getTotalCountByAuthorIDFunc func(ctx context.Context, authorID int) (int, error)
}

func (m *mockPostRepo) Create(ctx context.Context, post *model.Post) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, post)
	}
	return nil
}

func (m *mockPostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockPostRepo) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockPostRepo) GetTotalCount(ctx context.Context) (int, error) {
	if m.getTotalCountFunc != nil {
		return m.getTotalCountFunc(ctx)
	}
	return 0, nil
}

func (m *mockPostRepo) Exists(ctx context.Context, id int) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	return false, nil
}

func (m *mockPostRepo) GetByAuthorID(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
	if m.getByAuthorIDFunc != nil {
		return m.getByAuthorIDFunc(ctx, authorID, limit, offset)
	}
	return nil, nil
}

func (m *mockPostRepo) GetTotalCountByAuthorID(ctx context.Context, authorID int) (int, error) {
	if m.getTotalCountByAuthorIDFunc != nil {
		return m.getTotalCountByAuthorIDFunc(ctx, authorID)
	}
	return 0, nil
}

func TestPostService_Create_Success(t *testing.T) {
	mockPostRepo := &mockPostRepo{
		createFunc: func(ctx context.Context, post *model.Post) error {
			post.ID = 1
			post.CreatedAt = time.Now()
			return nil
		},
	}
	mockUserRepo := &mockUserRepo{} // Not used in create

	service := NewPostService(mockPostRepo, mockUserRepo)

	req := &model.PostCreateRequest{
		Title:   "Test Title",
		Content: "Test Content",
	}

	result, err := service.Create(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected post ID 1, got %d", result.ID)
	}

	if result.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %s", result.Title)
	}
}

func TestPostService_Create_InvalidRequest(t *testing.T) {
	mockPostRepo := &mockPostRepo{}
	mockUserRepo := &mockUserRepo{}

	service := NewPostService(mockPostRepo, mockUserRepo)

	req := &model.PostCreateRequest{
		Title:   "",
		Content: "Test Content",
	}

	_, err := service.Create(context.Background(), 1, req)
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestPostService_GetByID_Success(t *testing.T) {
	mockPost := &model.Post{
		ID:        1,
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
	}
	mockPostRepo := &mockPostRepo{
		getByIDFunc: func(ctx context.Context, id int) (*model.Post, error) {
			return mockPost, nil
		},
	}
	mockUserRepo := &mockUserRepo{}

	service := NewPostService(mockPostRepo, mockUserRepo)

	result, err := service.GetByID(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected post ID 1, got %d", result.ID)
	}
}

func TestPostService_GetByID_NotFound(t *testing.T) {
	mockPostRepo := &mockPostRepo{
		getByIDFunc: func(ctx context.Context, id int) (*model.Post, error) {
			return nil, errors.New("post not found")
		},
	}
	mockUserRepo := &mockUserRepo{}

	service := NewPostService(mockPostRepo, mockUserRepo)

	_, err := service.GetByID(context.Background(), 1, 1)
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestPostService_GetAll(t *testing.T) {
	mockPosts := []*model.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", AuthorID: 1},
		{ID: 2, Title: "Post 2", Content: "Content 2", AuthorID: 1},
	}
	mockPostRepo := &mockPostRepo{
		getAllFunc: func(ctx context.Context, limit, offset int) ([]*model.Post, error) {
			return mockPosts, nil
		},
		getTotalCountFunc: func(ctx context.Context) (int, error) {
			return 2, nil
		},
	}
	mockUserRepo := &mockUserRepo{}

	service := NewPostService(mockPostRepo, mockUserRepo)

	posts, total, err := service.GetAll(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestPostService_GetByAuthor(t *testing.T) {
	mockPosts := []*model.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", AuthorID: 1},
	}
	mockPostRepo := &mockPostRepo{
		getByAuthorIDFunc: func(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
			return mockPosts, nil
		},
		getTotalCountByAuthorIDFunc: func(ctx context.Context, authorID int) (int, error) {
			return 1, nil
		},
	}
	mockUserRepo := &mockUserRepo{}

	service := NewPostService(mockPostRepo, mockUserRepo)

	posts, total, err := service.GetByAuthor(context.Background(), 1, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}

	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
}
