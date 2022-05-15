package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-pg/pg/v10"

	"github.com/quantonganh/talkie"
)

type userService struct {
	DB *pg.DB
}

func NewUserService(db *pg.DB) *userService {
	return &userService{
		DB: db,
	}
}

// Create creates new user in the database
func (us *userService) Create(ctx context.Context, userReq *talkie.CreateUserRequest) (*talkie.User, error) {
	query, args, err := sq.
		Insert("user_account").
		Columns(
			"name",
			"email",
			"profile_picture",
			"provider",
			"provider_id",
		).
		Values(
			userReq.Name,
			userReq.Email,
			userReq.ProfilePicture,
			userReq.Provider,
			userReq.ProviderID,
		).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user talkie.User
	_, err = us.DB.QueryOneContext(ctx, &user, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}

// GetByID gets user by its id
// Returns ErrNotFound if id does not exist
func (us *userService) GetByID(ctx context.Context, id int) (*talkie.User, error) {
	query, args, err := sq.
		Select(
			"id",
			"name",
			"email",
			"profile_picture",
			"provider",
			"provider_id",
		).
		From("user_account").
		Where(sq.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user talkie.User
	_, err = us.DB.QueryOneContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, &talkie.Error{
				Err:  talkie.ErrUserNotFound,
				Code: talkie.ErrNotFound,
			}
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}

// GetByProviderID gets user by its provider id
// Returns ErrNotFound if provider id does not exist
func (us *userService) GetByProviderID(ctx context.Context, provider, providerID string) (*talkie.User, error) {
	query, args, err := sq.
		Select(
			"id",
			"name",
			"email",
			"profile_picture",
			"provider",
			"provider_id",
		).
		From("user_account").
		Where(sq.And{
			sq.Eq{
				"provider": provider,
			},
			sq.Eq{
				"provider_id": providerID,
			},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user talkie.User
	_, err = us.DB.QueryOneContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, &talkie.Error{
				Err:  talkie.ErrUserNotFound,
				Code: talkie.ErrNotFound,
			}
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}
