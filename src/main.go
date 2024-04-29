package main

import (
	"go-api/src/middleware"
	"go-api/src/routes"
	"go-api/src/utils"
	"log"
	"net/http"
)

type AppData struct {
}

func main() {
	//init the env variables
	env, envErr := utils.InitEnv()
	if envErr != nil {
		log.Panicln(envErr)
	}

	//create router
	mainRouter := http.NewServeMux()

	//register middlewares
	logger := middleware.Logging

	middlewaresStack := middleware.CreateStack(logger)
	middlewaresStack(mainRouter)

	//register routing
	routes.RegisterRoutes(mainRouter)

	//create server instance
	app := http.Server{
		Addr:    ":" + env.PORT,
		Handler: mainRouter,
	}

	//start listening for traffic
	log.Println("starting server on port " + env.PORT)
	listenErr := app.ListenAndServe()
	if listenErr != nil {
		log.Panicln(listenErr)
	}
}
