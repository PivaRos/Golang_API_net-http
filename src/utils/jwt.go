package utils

import (
	"fmt"
	"go-api/src/enums"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	UserId string
	Role   enums.Role
	jwt.StandardClaims
}

type AuthClaims struct {
	Otp    string
	UserId string
	Role   enums.Role
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func GetUserClaims(token string, secretKey []byte) (*UserClaims, error) {
	claims := &UserClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(parsedToken *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}

func GetAuthClaims(token string, secretKey []byte) (*AuthClaims, error) {
	claims := &AuthClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(parsedToken *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}
