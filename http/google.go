package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"

	"github.com/quantonganh/talkie"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type claims struct {
	UserID         int    `json:"user_id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
	jwt.StandardClaims
}

func (s *Server) authGoogle() appHandler {
	return func(c *gin.Context) error {
		var req talkie.AuthGoogleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return err
		}

		ctx := context.Background()
		payload, err := idtoken.Validate(ctx, req.Credential, os.Getenv("GOOGLE_CLIENT_ID"))
		if err != nil {
			return err
		}

		userReq := talkie.CreateUserRequest{
			Provider:   talkie.ProviderGoogle,
			ProviderID: payload.Subject,
		}
		username, ok := payload.Claims["name"]
		if ok {
			userReq.Name = username.(string)
		} else {
			return errors.New("name not found in claims")
		}

		email, ok := payload.Claims["email"]
		if ok {
			userReq.Email = email.(string)
		} else {
			return errors.New("email not found in claims")
		}

		picture, ok := payload.Claims["picture"]
		if ok {
			userReq.ProfilePicture = picture.(string)
		} else {
			return errors.New("picture not found in claims")
		}

		var user *talkie.User
		user, err = s.UserService.GetByProviderID(ctx, talkie.ProviderGoogle, payload.Subject)
		if err != nil {
			clientErr := new(talkie.Error)
			if errors.As(err, &clientErr) {
				if errors.Is(clientErr.Err, talkie.ErrUserNotFound) {
					user, err = s.UserService.Create(ctx, &userReq)
					if err != nil {
						return err
					}
				}
			}
			return err
		}

		token, err := createToken(user.ID, user.Email, user.Name, user.ProfilePicture)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, token)
		return nil
	}
}

func createToken(userID int, email, username, profilePicture string) (*Token, error) {
	atExpirationTime := time.Now().Add(15 * time.Minute)
	atClaims := &claims{
		UserID:         userID,
		Email:          email,
		Username:       username,
		ProfilePicture: profilePicture,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpirationTime.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return nil, err
	}

	rtExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	rtClaims := &claims{
		UserID:         userID,
		Email:          email,
		Username:       username,
		ProfilePicture: profilePicture,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
