package talkie

import (
	"context"
	"time"
)

const (
	ProviderGoogle = "google"
)

// User represents a social user
type User struct {
	ID         int       `json:"id" db:"id"`
	Provider   string    `json:"provider" db:"provider"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	UserComment
}

type UserComment struct {
	Email          string `json:"email" db:"email"`
	Name           string `json:"name" db:"name"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
}

// UserService is the interface that wraps the CRUD methods
type UserService interface {
	Create(ctx context.Context, user *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByProviderID(ctx context.Context, provider, providerID string) (*User, error)
}

// CreateUserRequest represents a request body when creating new user
type CreateUserRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
	Provider       string `json:"provider"`
	ProviderID     string `json:"provider_id"`
}
