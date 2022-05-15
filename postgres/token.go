package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-pg/pg/v10"

	"github.com/quantonganh/talkie"
)

type tokenService struct {
	DB *pg.DB
}

func NewTokenService(db *pg.DB) *tokenService {
	return &tokenService{
		DB: db,
	}
}

// Create creates new refresh token in the database
func (us *tokenService) Create(ctx context.Context, tokenReq *talkie.CreateTokenRequest) error {
	query, args, err := sq.
		Insert("token").
		Columns(
			"user_id",
			"refresh_token",
		).
		Values(
			tokenReq.UserID,
			tokenReq.RefreshToken,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = us.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
