package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-pg/pg/v10"

	"github.com/quantonganh/talkie"
)

const (
	ErrMsgParentIDInvalid = "parent_id=%d is not present in table"
)

type commentService struct {
	DB *pg.DB
}

func NewCommentService(db *pg.DB) *commentService {
	return &commentService{
		DB: db,
	}
}

// Create creates new comment in the database
// Returns ErrInvalid code if parent id is not present in the table
func (cs *commentService) Create(ctx context.Context, comment *talkie.CreateCommentRequest) error {
	query, args, err := sq.
		Insert("comment").
		Columns(
			"post_slug",
			"user_id",
			"content",
			"parent_id",
		).
		Values(
			comment.PostSlug,
			comment.UserID,
			comment.Content,
			comment.ParentID,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = cs.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr pg.Error
		if errors.As(err, &pgErr) {
			if pgErr.IntegrityViolation() {
				return &talkie.Error{
					Code:    talkie.ErrInvalid,
					Message: fmt.Sprintf(ErrMsgParentIDInvalid, *comment.ParentID),
				}
			}
		}
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// GetByPostSlug gets comment (and all sub-comments) by post slug
// Returns ErrNotFound if id does not exist
func (cs *commentService) GetByPostSlug(ctx context.Context, postSlug string) ([]*talkie.Comment, error) {
	query, args, err := sq.
		Select(
			"c.*",
			"comments(c.id) AS comments",
			"u.profile_picture",
			"u.name",
		).
		From("comment AS c").
		LeftJoin("user_account AS u ON u.id = c.user_id").
		Where(sq.And{
			sq.Eq{
				"post_slug": postSlug,
			},
			sq.Eq{
				"parent_id": nil,
			},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	comments := make([]*talkie.Comment, 0)
	_, err = cs.DB.QueryContext(ctx, &comments, query, args...)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, &talkie.Error{
				Err:  talkie.ErrCommentNotFound,
				Code: talkie.ErrNotFound,
			}
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return comments, nil
}

// GetByID gets comment (and all sub-comments) by its id
// Returns ErrNotFound if id does not exist
func (cs *commentService) GetByID(ctx context.Context, id int) (*talkie.Comment, error) {
	query, args, err := sq.
		Select(
			"id",
			"post_slug",
			"user_id",
			"content",
			"comments(id) AS comments",
		).
		From("comment").
		Where(sq.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var comment talkie.Comment
	_, err = cs.DB.QueryOneContext(ctx, &comment, query, args...)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, &talkie.Error{
				Err:  talkie.ErrCommentNotFound,
				Code: talkie.ErrNotFound,
			}
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &comment, nil
}

// UpdateByID updates a comment by id
func (cs *commentService) UpdateByID(ctx context.Context, id int, comment *talkie.Comment) error {
	query, args, err := sq.
		Update("comment").
		SetMap(map[string]interface{}{
			"content":    comment.Content,
			"updated_at": time.Now(),
		}).
		Where(sq.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = cs.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// List lists all comments
func (cs *commentService) List(ctx context.Context) ([]*talkie.Comment, error) {
	query, args, err := sq.
		Select(
			"id",
			"post_slug",
			"user_id",
			"created_at",
			"content",
			"comments(id) AS comments",
		).
		From("comment").
		Where(sq.Eq{
			"parent_id": nil,
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	comments := make([]*talkie.Comment, 0)
	_, err = cs.DB.QueryContext(ctx, &comments, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return comments, nil
}
