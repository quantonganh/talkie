package http

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/quantonganh/talkie"
)

func (s *Server) refreshTokenHandler() appHandler {
	return func(c *gin.Context) error {
		refreshTokenString := strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]

		rtClaims := &claims{}
		refreshToken, err := jwt.ParseWithClaims(refreshTokenString, rtClaims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return &talkie.Error{
					Code: talkie.ErrInvalid,
					Err:  jwt.ErrSignatureInvalid,
				}
			}
			return err
		}

		if !refreshToken.Valid {
			return &talkie.Error{
				Code: talkie.ErrUnauthorized,
			}
		}

		token, err := createToken(rtClaims.UserID, rtClaims.Email, rtClaims.Username, rtClaims.ProfilePicture)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, token)
		return nil
	}
}
