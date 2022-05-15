package talkie

import (
	"context"
	"time"
)

// Comment represents a post comment
type Comment struct {
	ID        int       `json:"id" db:"id"`
	PostSlug  string    `json:"post_slug" db:"post_slug"`
	UserID    int       `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	ParentID  *int      `json:"parent_id,omitempty" db:"parent_id"`
	Comments  []Comment `pg:",array" json:"comments,omitempty" db:"comments"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	UserComment `db:"user_account"`
}

// CommentService is the interface that wraps the CRUD methods
type CommentService interface {
	Create(ctx context.Context, comment *CreateCommentRequest) error
	GetByPostSlug(ctx context.Context, postSlug string) ([]*Comment, error)
	GetByID(ctx context.Context, id int) (*Comment, error)
	UpdateByID(ctx context.Context, id int, comment *Comment) error
	List(ctx context.Context) ([]*Comment, error)
}

// AuthGoogleRequest represents a request body when creating new comment
type AuthGoogleRequest struct {
	Credential string `json:"credential"`
}

// CreateCommentRequest represents a request body when creating new comment
type CreateCommentRequest struct {
	PostSlug string `json:"post_slug"`
	UserID   int    `json:"user_id"`
	Content  string `json:"content"`
	ParentID *int   `json:"parent_id,omitempty"`
}

// GetCommentPathParam represents path parameters when getting comments of a post
type GetCommentPathParam struct {
	PostSlug string `json:"post_slug" uri:"post_slug"`
}

// CommentPathParam represents path parameters when editing a comment
type CommentPathParam struct {
	ID int `json:"id" uri:"id"`
}

// EditCommentRequest represents a request body when editing a comment
type EditCommentRequest struct {
	Content string `json:"content"`
}
