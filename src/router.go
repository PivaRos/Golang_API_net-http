package main

import (
	"go-api/src/routes"
	"go-api/src/utils"
	"net/http"
)

func loadRoutes(router *http.ServeMux, appData *utils.AppData) {

	//auth
	routes.RegisterAuthRoutes(router, appData)
	//admin
	routes.RegisterAdminRoutes(router, appData)
	//appointment
	routes.RegisterAppointmentRoutes(router, appData)
}
