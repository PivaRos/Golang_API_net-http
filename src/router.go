package main

import (
	"go-api/src/auth"
	"go-api/src/care"
	"go-api/src/enums"
	"go-api/src/middleware"
	"go-api/src/utils"
	"net/http"
)

func loadRoutes(router *http.ServeMux, appData *utils.AppData) {
	authRouter := http.NewServeMux()
	adminRouter := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	var adminAccess []enums.Role = []enums.Role{
		enums.Admin,
	}
	router.Handle("/admin/", http.StripPrefix("/admin", middleware.Authenticate(adminAccess, appData)(adminRouter)))

	auth := auth.CreateHandler(appData)
	authRouter.HandleFunc("POST /Login", auth.Login)

	care := care.CreateHandler(appData.Database)
	adminRouter.HandleFunc("GET /Care", care.Get)
	adminRouter.HandleFunc("GET /Care/{id}", care.Get)
	adminRouter.HandleFunc("POST /Care", care.Create)
	adminRouter.HandleFunc("DELETE /Care/{id}", care.Delete)

}
