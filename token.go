package talkie

import (
	"context"
	"time"
)

// Token represents a refresh token
type Token struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// TokenService is the interface that wraps the CRUD methods
type TokenService interface {
	Create(ctx context.Context, user *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByProviderID(ctx context.Context, provider, providerID string) (*User, error)
}

// CreateTokenRequest represents a request body when creating new token
type CreateTokenRequest struct {
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest represents a request body when rotating tokens
type RefreshTokenRequest struct {
	PostSlug     string `json:"post_slug"`
	UserID       string `json:"user_id"`
	Content      string `json:"content"`
	ParentID     *int   `json:"parent_id,omitempty"`
	RefreshToken string `json:"refresh_token"`
}
