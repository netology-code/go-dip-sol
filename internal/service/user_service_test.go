package service

import (
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/pkg/auth"
	"context"
	"errors"
	"testing"
	"time"
)

// mockUserRepo is a mock implementation of UserRepository
type mockUserRepo struct {
	createFunc            func(ctx context.Context, user *model.User) error
	getByIDFunc           func(ctx context.Context, id int) (*model.User, error)
	getByEmailFunc        func(ctx context.Context, email string) (*model.User, error)
	getByUsernameFunc     func(ctx context.Context, username string) (*model.User, error)
	existsByEmailFunc     func(ctx context.Context, email string) (bool, error)
	existsByUsernameFunc  func(ctx context.Context, username string) (bool, error)
	updateFunc            func(ctx context.Context, user *model.User) error
	deleteFunc            func(ctx context.Context, id int) error
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.getByUsernameFunc != nil {
		return m.getByUsernameFunc(ctx, username)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.existsByEmailFunc != nil {
		return m.existsByEmailFunc(ctx, email)
	}
	return false, nil
}

func (m *mockUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if m.existsByUsernameFunc != nil {
		return m.existsByUsernameFunc(ctx, username)
	}
	return false, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *model.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}



func TestUserService_Register_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		existsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		existsByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, user *model.User) error {
			user.ID = 1
			user.CreatedAt = time.Now()
			user.UpdatedAt = time.Now()
			return nil
		},
	}
	jwtManager := auth.NewJWTManager("test-secret", 24)

	service := NewUserService(mockRepo, jwtManager)

	req := &model.UserCreateRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	result, err := service.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Token == "" {
		t.Error("expected non-empty token")
	}

	if result.User.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", result.User.Username)
	}
}

func TestUserService_Register_EmailExists(t *testing.T) {
	mockRepo := &mockUserRepo{
		existsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
	}
	jwtManager := auth.NewJWTManager("test-secret", 24)

	service := NewUserService(mockRepo, jwtManager)

	req := &model.UserCreateRequest{
		Username: "testuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	_, err := service.Register(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserService_Login_Success(t *testing.T) {
	hashedPassword, _ := auth.HashPassword("password123")
	mockRepo := &mockUserRepo{
		getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:        1,
				Username:  "testuser",
				Email:     email,
				Password:  hashedPassword,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	jwtManager := auth.NewJWTManager("test-secret", 24)

	service := NewUserService(mockRepo, jwtManager)

	req := &model.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	result, err := service.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	hashedPassword, _ := auth.HashPassword("password123")
	mockRepo := &mockUserRepo{
		getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:        1,
				Username:  "testuser",
				Email:     email,
				Password:  hashedPassword,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	jwtManager := auth.NewJWTManager("test-secret", 24)

	service := NewUserService(mockRepo, jwtManager)

	req := &model.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := service.Login(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
