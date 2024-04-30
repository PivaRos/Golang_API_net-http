package main

import (
	"go-api/src/auth"
	"go-api/src/care"
	"go-api/src/middleware"
	"net/http"
)

func loadRoutes(router *http.ServeMux, appData *AppData) {
	authRouter := http.NewServeMux()
	adminRouter := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	router.Handle("/admin/", http.StripPrefix("/admin", middleware.CheckAdmin(adminRouter)))

	auth := auth.CreateHandler(appData.Database)
	authRouter.HandleFunc("POST /auth/Login", auth.Login)

	care := care.CreateHandler(appData.Database)
	adminRouter.HandleFunc("POST /Care", care.Create)

}
