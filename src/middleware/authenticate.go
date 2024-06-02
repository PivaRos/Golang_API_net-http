package middleware

import (
	"context"
	"encoding/json"
	"go-api/src/enums"
	"go-api/src/user"
	"go-api/src/utils"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

func Authenticate(roles []enums.Role, app *utils.AppData) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &utils.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				return app.Env.Jwt_Secret_Key, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claim, ok := token.Claims.(*utils.UserClaims); ok && token.Valid {
				var found bool = false
				for _, value := range roles {
					if value == claim.Role {
						found = true
					}
				}
				if found {
					filter := bson.M{"accessToken": tokenString}
					var user user.User
					err := app.MongoClient.Database(app.Env.Db).Collection("users").FindOne(context.TODO(), filter).Decode(&user)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					bytes, _ := json.Marshal(user)
					log.Println("user1", string(bytes))
					err = user.Validate()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					ctx := context.WithValue(r.Context(), utils.UserDataContextKey, user)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		})
	}
}
