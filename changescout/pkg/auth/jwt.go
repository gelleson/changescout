package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTWithClaims[T jwt.Claims](claims T, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
