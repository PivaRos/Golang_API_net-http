package main

import (
	"go-api/src/middleware"
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

	//register routing
	loadRoutes(mainRouter)

	//register middlewares
	logger := middleware.Logging
	Stack := middleware.CreateStack(logger)

	//create server instance
	app := http.Server{
		Addr:    ":" + env.PORT,
		Handler: Stack(mainRouter),
	}

	//start listening for traffic
	log.Println("starting server on port " + env.PORT)
	listenErr := app.ListenAndServe()
	if listenErr != nil {
		log.Panicln(listenErr)
	}
}
