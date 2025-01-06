package middleware

import (
	"errors"
	"golang-gorm/domain/model"
	"golang-gorm/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type authMiddleware struct {
	secretKeyUser string
}

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddleware{
		secretKeyUser: viper.GetString("JWT_SECRET_KEY_USER"),
	}
}

type AuthMiddleware interface {
	AuthUser() gin.HandlerFunc
}

func (m *authMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		requestToken := c.Request.Header.Get("Authorization")
		if requestToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: Missing Authorization header",
			})
			return
		}

		// check token format
		splitToken := strings.Split(requestToken, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: Invalid token format",
			})
			return
		}

		// get token without 'Bearer '
		tokenString := splitToken[1]

		// Validate token
		token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeyUser), nil
		})

		// check validity token
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.Response{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized: Invalid token signature",
				})
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.Response{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized: Token expired",
				})
				return
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.Response{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		claims, ok := token.Claims.(*model.JWTClaimUser)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized: Invalid token claims",
			})
			return
		}

		c.Set("user_data", *claims)
		c.Next()
	}
}
