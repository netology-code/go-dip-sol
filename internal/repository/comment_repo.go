package repository

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (post_id, author_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		comment.PostID, comment.AuthorID, comment.Content,
		comment.CreatedAt, comment.UpdatedAt,
	).Scan(&comment.ID)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	query := `
		SELECT id, post_id, author_id, content, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	var comment model.Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments: %w", err)
	}

	return comments, nil
}

func (r *CommentRepo) GetCountByPostID(ctx context.Context, postID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}
	return count, nil
}
