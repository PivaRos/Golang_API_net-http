package main

import (
	"context"
	"go-api/src/middleware"
	"go-api/src/utils"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
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
	//setup redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.Redis_Addr,
		Password: env.Redis_Password,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("Connected to redis")
			return nil
		},
		DB: 0,
	})
	rError := utils.CheckRedisConnection(context.TODO(), rdb)
	if rError != nil {
		log.Panicln(err)
	}
	appData := &utils.AppData{
		MongoClient: client,
		Database:    db,
		Env:         env,
		RedisClient: rdb,
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
