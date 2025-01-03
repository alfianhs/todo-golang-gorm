package model

import "github.com/golang-jwt/jwt/v5"

type JWTClaimUser struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
