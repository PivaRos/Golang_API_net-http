package main

import (
	"go-api/src/middleware"
	"go-api/src/utils"
	"log"
	"net/http"
)

func main() {
	//init the env variables
	env, err := utils.InitEnv()
	if err != nil {
		log.Panicln(err)
	}

	//setup mongodb connection
	client, db, currentContext, cancel, err := utils.SetupMongoDB(env.MONGO_URI, env.Db)
	defer utils.CloseConnection(client, currentContext, cancel)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Connected to mongodb")

	appData := &utils.AppData{
		MongoClient: client,
		Database:    db,
		Env:         env,
	}

	//create router
	mainRouter := http.NewServeMux()

	//register routing
	loadRoutes(mainRouter, appData)

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
	err = app.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
}
