package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/quantonganh/talkie"
)

var (
	ErrMarshalFailed = "failed to marshal: %w"
	ErrDecodeFailed  = "failed to decode: %w"
)

func (s *Server) createCommentHandler() appHandler {
	return func(c *gin.Context) error {
		var req talkie.CreateCommentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return err
		}

		tokenString := strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]

		claims := &claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil {
			var vErr *jwt.ValidationError
			if errors.As(err, &vErr) {
				switch vErr.Errors {
				case jwt.ValidationErrorSignatureInvalid:
					return &talkie.Error{
						Code:    talkie.ErrInvalid,
						Err:     vErr,
						Message: "invalid signature",
					}
				case jwt.ValidationErrorExpired:
					return &talkie.Error{
						Code: talkie.ErrUnauthorized,
						Err:  vErr,
					}
				}
			}
			return err
		}

		if !token.Valid {
			return &talkie.Error{
				Code:    talkie.ErrInvalid,
				Message: "invalid token",
			}
		}

		req.UserID = claims.UserID
		if err := s.CommentService.Create(context.Background(), &req); err != nil {
			return err
		}

		c.JSON(http.StatusOK, req)
		return nil
	}
}

func (s *Server) getCommentsHandler() appHandler {
	return func(c *gin.Context) error {
		var param talkie.GetCommentPathParam
		if err := c.BindUri(&param); err != nil {
			return err
		}

		comment, err := s.CommentService.GetByPostSlug(context.Background(), param.PostSlug)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, comment)
		return nil
	}
}

func (s *Server) updateCommentHandler() appHandler {
	return func(c *gin.Context) error {
		var param talkie.CommentPathParam
		if err := c.BindUri(&param); err != nil {
			return err
		}

		var req talkie.EditCommentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return err
		}
		body, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf(ErrMarshalFailed, err)
		}

		ctx := context.Background()
		comment, err := s.CommentService.GetByID(ctx, param.ID)
		if err != nil {
			return err
		}

		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&comment); err != nil {
			return fmt.Errorf(ErrDecodeFailed, err)
		}

		err = s.CommentService.UpdateByID(ctx, param.ID, comment)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, comment)
		return nil
	}
}

func (s *Server) listCommentsHandler() appHandler {
	return func(c *gin.Context) error {
		comments, err := s.CommentService.List(context.Background())
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, comments)
		return nil
	}
}
