package helpers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GetJWTSecretKeyUser() string {
	return viper.GetString("JWT_SECRET_KEY_USER")
}

func GetJWTTTL() int {
	return viper.GetInt("JWT_TTL")
}

func GenerateJWTTokenUser(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeyUser()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
