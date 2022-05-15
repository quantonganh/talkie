package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func authJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
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
					c.AbortWithStatus(http.StatusBadRequest)
				case jwt.ValidationErrorExpired:
					c.AbortWithStatus(http.StatusUnauthorized)
				}
			}
		}

		if !token.Valid {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		c.Next()
	}
}
