package main

import (
	"go-api/src/appointment"
	"go-api/src/auth"
	"go-api/src/care"
	"go-api/src/enums"
	"go-api/src/middleware"
	"go-api/src/utils"
	"net/http"
)

func loadRoutes(router *http.ServeMux, appData *utils.AppData) {

	//auth
	authRouter := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	auth := auth.CreateHandler(appData)
	authRouter.HandleFunc("POST /", auth.SendOTP)
	authRouter.HandleFunc("POST /Validate", auth.ValidateOTP)

	//admin
	adminRouter := http.NewServeMux()
	var adminAccess []enums.Role = []enums.Role{
		enums.Admin,
	}
	router.Handle("/admin/", http.StripPrefix("/admin", middleware.Authenticate(adminAccess, appData)(adminRouter)))
	//admin -> care
	care := care.CreateHandler(appData.Database)
	adminRouter.HandleFunc("GET /care", care.Get)
	adminRouter.HandleFunc("GET /care/{id}", care.Get)
	adminRouter.HandleFunc("POST /care", care.Create)
	adminRouter.HandleFunc("DELETE /care/{id}", care.Delete)

	//appointment
	appointmentRouter := http.NewServeMux()
	var allAccess []enums.Role = []enums.Role{
		enums.Customer,
	}
	router.Handle("/appointment/", http.StripPrefix("/appointment", middleware.Authenticate(allAccess, appData)(appointmentRouter)))
	appointment := appointment.CreateHandler(appData.Database, *appData.RedisClient)
	appointmentRouter.HandleFunc("GET /", appointment.Get)
	appointmentRouter.HandleFunc("GET /times", appointment.GetAvailableTime)
	appointmentRouter.HandleFunc("GET /{id}", appointment.GetById)

}
