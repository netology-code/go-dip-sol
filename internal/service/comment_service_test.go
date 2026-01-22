package service

import (
	"advanced-blog-management-system/internal/model"

	"context"
	"errors"
	"testing"
	"time"
)

// mockCommentRepo is a mock implementation of CommentRepository
type mockCommentRepo struct {
	createFunc             func(ctx context.Context, comment *model.Comment) error
	getByIDFunc            func(ctx context.Context, id int) (*model.Comment, error)
	getByPostIDFunc        func(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error)
	getCountByPostIDFunc   func(ctx context.Context, postID int) (int, error)
}

func (m *mockCommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, comment)
	}
	return nil
}

func (m *mockCommentRepo) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
	if m.getByPostIDFunc != nil {
		return m.getByPostIDFunc(ctx, postID, limit, offset)
	}
	return nil, nil
}

func (m *mockCommentRepo) GetCountByPostID(ctx context.Context, postID int) (int, error) {
	if m.getCountByPostIDFunc != nil {
		return m.getCountByPostIDFunc(ctx, postID)
	}
	return 0, nil
}

func TestCommentService_Create_Success(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{
		createFunc: func(ctx context.Context, comment *model.Comment) error {
			comment.ID = 1
			comment.CreatedAt = time.Now()
			comment.UpdatedAt = time.Now()
			return nil
		},
	}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
	}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	result, err := service.Create(context.Background(), 1, 1, "Test comment content")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected comment ID 1, got %d", result.ID)
	}

	if result.Content != "Test comment content" {
		t.Errorf("expected content 'Test comment content', got %s", result.Content)
	}
}

func TestCommentService_Create_InvalidPostID(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{}
	mockPostRepo := &mockPostRepo{}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	_, err := service.Create(context.Background(), 1, 0, "Test comment")
	if err == nil {
		t.Fatal("expected error for invalid post ID, got nil")
	}
}

func TestCommentService_Create_PostNotFound(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return false, nil
		},
	}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	_, err := service.Create(context.Background(), 1, 1, "Test comment")
	if err == nil {
		t.Fatal("expected error for post not found, got nil")
	}
}

func TestCommentService_Create_EmptyContent(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
	}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	_, err := service.Create(context.Background(), 1, 1, "   ")
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}
}

func TestCommentService_Create_ContentTooLong(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
	}

	longContent := string(make([]byte, 1001))
	service := NewCommentService(mockCommentRepo, mockPostRepo)

	_, err := service.Create(context.Background(), 1, 1, longContent)
	if err == nil {
		t.Fatal("expected error for content too long, got nil")
	}
}

func TestCommentService_GetByPost_Success(t *testing.T) {
	mockComments := []*model.Comment{
		{ID: 1, Content: "Comment 1", PostID: 1, AuthorID: 1},
		{ID: 2, Content: "Comment 2", PostID: 1, AuthorID: 2},
	}
	mockCommentRepo := &mockCommentRepo{
		getByPostIDFunc: func(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
			return mockComments, nil
		},
		getCountByPostIDFunc: func(ctx context.Context, postID int) (int, error) {
			return 2, nil
		},
	}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
	}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	comments, total, err := service.GetByPost(context.Background(), 1, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestCommentService_GetByPost_InvalidPostID(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{}
	mockPostRepo := &mockPostRepo{}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	_, _, err := service.GetByPost(context.Background(), 0, 10, 0)
	if err == nil {
		t.Fatal("expected error for invalid post ID, got nil")
	}
}

func TestCommentService_GetByPost_PostNotFound(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return false, nil
		},
	}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	_, _, err := service.GetByPost(context.Background(), 1, 10, 0)
	if err == nil {
		t.Fatal("expected error for post not found, got nil")
	}
}

func TestCommentService_GetByPost_LimitBounds(t *testing.T) {
	mockCommentRepo := &mockCommentRepo{
		getByPostIDFunc: func(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
			// Check that limit is bounded
			if limit < 1 || limit > 100 {
				t.Errorf("limit should be between 1 and 100, got %d", limit)
			}
			if offset < 0 {
				t.Errorf("offset should be >= 0, got %d", offset)
			}
			return []*model.Comment{}, nil
		},
		getCountByPostIDFunc: func(ctx context.Context, postID int) (int, error) {
			return 0, nil
		},
	}
	mockPostRepo := &mockPostRepo{
		existsFunc: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
	}

	service := NewCommentService(mockCommentRepo, mockPostRepo)

	// Test with limit < 1
	_, _, err := service.GetByPost(context.Background(), 1, 0, -1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test with limit > 100
	_, _, err = service.GetByPost(context.Background(), 1, 200, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
