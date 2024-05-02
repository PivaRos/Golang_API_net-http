package utils

import (
	"go-api/src/role"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   role.Role
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
