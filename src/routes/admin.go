package routes

import (
	"go-api/src/appointment"
	"go-api/src/care"
	"go-api/src/enums"
	"go-api/src/middleware"
	"go-api/src/utils"
	"net/http"
)

func RegisterAdminRoutes(router *http.ServeMux, appData *utils.AppData) {
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

	appointment := appointment.CreateHandler(appData.Database)
	adminRouter.HandleFunc("GET /appointment", appointment.AdminGet)
	adminRouter.HandleFunc("GET /appointment/times", appointment.GetAvailableTime)
	adminRouter.HandleFunc("GET /appointment/{id}", appointment.AdminGetById)
}
