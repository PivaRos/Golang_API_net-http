package auth

import (
	"context"
	"errors"
	"go-api/src/enums"
	"go-api/src/user"
	"go-api/src/utils"
	"log"
	"math/rand"
	"strconv"
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

func (s *services) SendOTP(l Login) (*string, error) {
	var user user.User
	collection := s.app.Database.Collection("Users")
	err := collection.FindOne(context.TODO(), bson.M{"phone": l.Phone, "govId": l.GovId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	code := strconv.Itoa(rand.Intn(900000) + 100000)
	if s.app.Env.ENVIRONMENT == "Development" {
		code = "111111" //! only for development
	}
	err = s.SendSMS(l.Phone, "", "Your Verification code is: "+code)
	if err != nil {
		return nil, err
	}
	authToken, err := s.GenerateAuthToken(user, code)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (s *services) ValidateOTP(token string, code string) (*utils.Tokens, error) {
	log.Println("jere-1")
	log.Println("token is : " + token)
	claims, err := utils.GetAuthClaims(token, s.app.Env.Jwt_Secret_Key)
	if err != nil {
		return nil, err
	}
	log.Println(claims.Otp)
	log.Println(code)
	if claims.Otp == code {
		// Generate new access and refresh tokens
		newTokens, err := s.GenerateUserTokens(claims.UserId, claims.Role)
		if err != nil {
			return nil, err
		}
		return &newTokens, nil
	}
	return nil, errors.New("otp is not matching")
}

func (s *services) GenerateUserTokens(userId string, role enums.Role) (utils.Tokens, error) {
	var tokens utils.Tokens
	// Set expiration times for each token
	accessTokenExpTime := s.app.Env.Access_Token_Expiration
	refreshTokenExpTime := s.app.Env.Refresh_Token_Expiration

	// Generate Access Token
	accessTokenClaims := &utils.UserClaims{
		Role:   role,
		UserId: userId,
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
	refreshTokenClaims := &utils.UserClaims{

		Role:   role,
		UserId: userId,
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

	Id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return tokens, err
	}
	result, err := s.app.Database.Collection("Users").UpdateByID(context.TODO(), Id, bson.M{"$set": bson.M{"accessToken": tokens.AccessToken, "refreshToken": tokens.RefreshToken}})
	if err != nil {
		return tokens, err
	}
	_ = result // just to remove message

	return tokens, nil
}

func (s *services) GenerateAuthToken(user user.User, Otp string) (*string, error) {
	accessTokenExpTime := s.app.Env.Access_Token_Expiration

	accessTokenClaims := &utils.AuthClaims{
		Otp:    Otp,
		UserId: user.Id.Hex(),
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	authToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(s.app.Env.Jwt_Secret_Key)
	if err != nil {
		return nil, err
	}
	return &authToken, nil
}

func (s *services) RefreshToken(oldRefreshToken string) (*utils.Tokens, error) {
	claims, err := utils.GetUserClaims(oldRefreshToken, s.app.Env.Jwt_Secret_Key)
	if err != nil {
		return nil, err
	}

	// Generate new access and refresh tokens
	newTokens, err := s.GenerateUserTokens(claims.UserId, claims.Role)
	if err != nil {
		return nil, err
	}

	return &newTokens, nil

}

func (s *services) InvalidateToken(token string) error {
	tokensCollection := s.app.Database.Collection("tokens")
	filter := bson.M{"token": token}
	result, err := tokensCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no documents found")
	}
	return nil
}

func (s *services) SendSMS(phone string, title string, body string) error {
	return nil
}
