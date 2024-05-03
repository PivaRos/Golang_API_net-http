package auth

import (
	"context"
	"fmt"
	"go-api/src/role"
	"go-api/src/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateServices(app *utils.AppData) *services {
	return &services{
		app: app,
	}
}

type services struct {
	app *utils.AppData
}

func (s *services) Login(l Login) (*utils.Tokens, error) {
	var user User
	collection := s.app.Database.Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"phone": l.Phone, "password": l.Password}).Decode(&user)
	if err != nil {
		return nil, err
	}
	tokens, err := s.GenerateTokens(user.Id.Hex(), user.Role)
	if err != nil {
		return nil, err
	}
	return &tokens, nil

}

func (s *services) GenerateTokens(userID string, role role.Role) (utils.Tokens, error) {
	var tokens utils.Tokens
	// Set expiration times for each token
	accessTokenExpTime := s.app.Env.Access_Token_Expiration
	refreshTokenExpTime := s.app.Env.Refresh_Token_Expiration

	// Generate Access Token
	accessTokenClaims := &utils.Claims{
		Role:   role,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(s.app.Env.Jwt_Secret_Key)
	if err != nil {
		return tokens, err
	}
	tokens.AccessToken = accessToken

	// Generate Refresh Token
	refreshTokenClaims := &utils.Claims{
		Role:   role,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenExpTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString(s.app.Env.Jwt_Secret_Key)
	if err != nil {
		return tokens, err
	}
	tokens.RefreshToken = refreshToken

	Id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return tokens, err
	}
	result, err := s.app.Database.Collection("users").UpdateByID(context.TODO(), Id, bson.M{"$set": bson.M{"accessToken": tokens.AccessToken, "refreshToken": tokens.RefreshToken}})
	if err != nil {
		return tokens, err
	}
	_ = result // just to remove message

	return tokens, nil
}

func (s *services) RefreshToken(oldRefreshToken string, role role.Role) (utils.Tokens, error) {
	var tokens utils.Tokens

	token, err := jwt.ParseWithClaims(oldRefreshToken, utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.app.Env.Jwt_Secret_Key, nil
	})
	if err != nil {
		return tokens, err
	}

	claims, ok := token.Claims.(utils.Claims)
	if !ok || !token.Valid {
		return tokens, fmt.Errorf("Invalid refresh token")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return tokens, fmt.Errorf("Expired refresh token")
	}

	// Generate new access and refresh tokens
	newTokens, err := s.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		return tokens, err
	}

	return newTokens, nil
}
