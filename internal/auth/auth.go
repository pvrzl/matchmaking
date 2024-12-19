package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -source=types.go -destination=./../mock/auth/types_mock.go -package=auth
type TokenGenerator interface {
	GenerateToken(userID int, email string) (string, error)
}

type JWTTokenGenerator struct {
	secretKey string
}

func NewJWTTokenGenerator(secretKey string) *JWTTokenGenerator {
	return &JWTTokenGenerator{secretKey: secretKey}
}

func (j *JWTTokenGenerator) GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
