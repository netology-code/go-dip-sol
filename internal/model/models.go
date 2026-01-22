package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	PostID    int       `json:"post_id" db:"post_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PostCreateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
	PostID  int    `json:"post_id" validate:"required,gt=0"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type PostResponse struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

type CommentResponse struct {
	ID        int          `json:"id"`
	Content   string       `json:"content"`
	PostID    int          `json:"post_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func (r *UserCreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UserLoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *PostCreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *CommentCreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
