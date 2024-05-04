package utils

import (
	"go-api/src/enums"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   enums.Role
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
