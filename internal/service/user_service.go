package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"advanced-blog-management-system/pkg/auth"
	"context"
	"errors"
	"fmt"
)

type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// 1. Валидация входных данных
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. Проверка уникальности email
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// 3. Проверка уникальности username
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// 4. Хеширование пароля
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Создание модели пользователя
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// 6. Сохранение пользователя
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 7. Генерация JWT токена
	token, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 8. Возврат TokenResponse
	return &model.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	// 1. Валидация входных данных
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. Поиск пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 3. Проверка пароля
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 4. Генерация JWT токена
	token, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 5. Возврат TokenResponse
	return &model.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
