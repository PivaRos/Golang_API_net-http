package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	MONGO_URI                string
	Jwt_Secret_Key           []byte
	Access_Token_Expiration  time.Duration
	Refresh_Token_Expiration time.Duration
	Db                       string
	PORT                     string
}

func InitEnv() (*Env, error) {

	e := Env{}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information")
	}
	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, "../../.env")
	if os.Getenv("EnvFile") == "" {
		EnvErr := godotenv.Load(envPath)
		if EnvErr != nil {
			return nil, EnvErr
		}
	}
	e.MONGO_URI = os.Getenv("MONGO_URI")
	if e.MONGO_URI == "" {
		return nil, errors.New("no uri was found in env file")
	}
	Jwt_Secret_Key := os.Getenv("JWT_SECRET_KEY")
	if Jwt_Secret_Key == "" {
		return nil, errors.New("no JWT_SECRET_KEY was found in env file")
	}
	e.Jwt_Secret_Key = []byte(Jwt_Secret_Key)
	Access_Token_Expiration, AccessExpirationErr := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRATION"))
	if AccessExpirationErr != nil || Access_Token_Expiration == time.Duration(0) {
		return nil, errors.New("no ACCESS_TOKEN_EXPIRATION was found in env file")
	}
	e.Access_Token_Expiration = Access_Token_Expiration
	Refresh_Token_Expiration, RefreshExpirationErr := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRATION"))
	if RefreshExpirationErr != nil || Refresh_Token_Expiration == time.Duration(0) {
		return nil, errors.New("no REFRESH_TOKEN_EXPIRATION was found in env file")
	}
	e.Refresh_Token_Expiration = Refresh_Token_Expiration
	e.Db = os.Getenv("DB")
	if e.Db == "" {
		return nil, errors.New("no DB was found in env file")
	}
	e.PORT = os.Getenv("PORT")
	if e.PORT == "" {
		return nil, errors.New("no PORT was found in env file")
	}

	return &e, nil
}
