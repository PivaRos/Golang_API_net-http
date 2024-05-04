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

func GetJwtClaims(token string, Jwt_Secret_Key string) (interface{}, error) {
	parsedToken, err := jwt.ParseWithClaims(token, jwt.StandardClaims{}, func(parsedToken *jwt.Token) (interface{}, error) {
		return Jwt_Secret_Key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(UserClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}
