package model

import (
	"testing"
)

func TestUserCreateRequest_Validate_Success(t *testing.T) {
	req := UserCreateRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserCreateRequest_Validate_InvalidUsername(t *testing.T) {
	req := UserCreateRequest{
		Username: "", // too short
		Email:    "test@example.com",
		Password: "password123",
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for empty username")
	}
}

func TestUserCreateRequest_Validate_InvalidEmail(t *testing.T) {
	req := UserCreateRequest{
		Username: "testuser",
		Email:    "invalid-email",
		Password: "password123",
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for invalid email")
	}
}

func TestUserCreateRequest_Validate_InvalidPassword(t *testing.T) {
	req := UserCreateRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "123", // too short
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for short password")
	}
}

func TestUserLoginRequest_Validate_Success(t *testing.T) {
	req := UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserLoginRequest_Validate_InvalidEmail(t *testing.T) {
	req := UserLoginRequest{
		Email:    "",
		Password: "password123",
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for empty email")
	}
}

func TestPostCreateRequest_Validate_Success(t *testing.T) {
	req := PostCreateRequest{
		Title:   "Test Title",
		Content: "Test content",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestPostCreateRequest_Validate_InvalidTitle(t *testing.T) {
	req := PostCreateRequest{
		Title:   "",
		Content: "Test content",
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for empty title")
	}
}

func TestPostCreateRequest_Validate_InvalidContent(t *testing.T) {
	req := PostCreateRequest{
		Title:   "Test Title",
		Content: "",
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for empty content")
	}
}

func TestCommentCreateRequest_Validate_Success(t *testing.T) {
	req := CommentCreateRequest{
		Content: "Test comment",
		PostID:  1,
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCommentCreateRequest_Validate_InvalidContent(t *testing.T) {
	req := CommentCreateRequest{
		Content: "",
		PostID:  1,
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for empty content")
	}
}

func TestCommentCreateRequest_Validate_InvalidPostID(t *testing.T) {
	req := CommentCreateRequest{
		Content: "Test comment",
		PostID:  0, // invalid
	}

	err := req.Validate()
	if err == nil {
		t.Error("expected validation error for invalid post ID")
	}
}
