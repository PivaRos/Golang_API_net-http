package middleware

import (
	"context"
	"go-api/src/role"
	"go-api/src/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

func Authenticate(roles []role.Role, app *utils.AppData) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return app.Env.Jwt_Secret_Key, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claim, ok := token.Claims.(*utils.Claims); ok && token.Valid {
				var found bool = false
				for _, value := range roles {
					if value == claim.Role {
						found = true
					}
				}
				if found {
					filter := bson.M{"accessToken": tokenString}
					result := app.MongoClient.Database(app.Env.Db).Collection("users").FindOne(context.TODO(), filter)
					if result.Err() != nil {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

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
