package repository

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// UserRepo представляет репозиторий для работы с пользователями
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo создает новый репозиторий пользователей
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create создает нового пользователя
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrUserNotFound
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *UserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetByUsername получает пользователя по username
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// ExistsByEmail проверяет существование пользователя по email
func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// ExistsByUsername проверяет существование пользователя по username
func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}

// Update обновляет данные пользователя
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// Delete удаляет пользователя
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}
