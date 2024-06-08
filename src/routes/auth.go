package routes

import (
	"go-api/src/auth"
	"go-api/src/utils"
	"net/http"
)

func RegisterAuthRoutes(router *http.ServeMux, appData *utils.AppData) {
	authRouter := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	auth := auth.CreateHandler(appData)
	authRouter.HandleFunc("POST /", auth.SendOTP)
	authRouter.HandleFunc("POST /Validate", auth.ValidateOTP)
}
