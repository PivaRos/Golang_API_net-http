package middleware

import (
	"log"
	"net/http"
)

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("this is check admin")
		next.ServeHTTP(w, r)

	})
}
