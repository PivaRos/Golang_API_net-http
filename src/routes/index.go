package routes

import (
	"net/http"
)

func RegisterRoutes(mainRoute *http.ServeMux) {

	auth := &Auth{}
	mainRoute.HandleFunc("POST /auth/login", auth.Login)
}
