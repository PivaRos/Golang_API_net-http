package routes

import (
	"go-api/src/appointment"
	"go-api/src/enums"
	"go-api/src/middleware"
	"go-api/src/utils"
	"net/http"
)

func RegisterAppointmentRoutes(router *http.ServeMux, appData *utils.AppData) {
	appointmentRouter := http.NewServeMux()
	var allAccess []enums.Role = []enums.Role{
		enums.Customer,
	}
	router.Handle("/appointment/", http.StripPrefix("/appointment", middleware.Authenticate(allAccess, appData)(appointmentRouter)))
	appointment := appointment.CreateHandler(appData.Database)
	appointmentRouter.HandleFunc("POST /a", appointment.Create)
	appointmentRouter.HandleFunc("GET /times", appointment.GetAvailableTime)
	appointmentRouter.HandleFunc("GET /", appointment.Get)
	appointmentRouter.HandleFunc("GET /{id}", appointment.GetById)
}
