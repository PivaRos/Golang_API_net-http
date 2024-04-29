package main

import (
	"go-api/src/admin"
	"go-api/src/auth"
	"go-api/src/middleware"
	"net/http"
)

func loadRoutes(router *http.ServeMux) {

	auth := &auth.Handler{}
	router.HandleFunc("POST /auth/Login", auth.Login)

	admin := &admin.Handler{}
	subRouter := http.NewServeMux()
	router.Handle("/admin/", http.StripPrefix("/admin", middleware.CheckAdmin(subRouter)))
	subRouter.HandleFunc("POST /CreateCare", admin.CreateCare)

}
