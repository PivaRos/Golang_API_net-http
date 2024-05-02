package utils

import "go.mongodb.org/mongo-driver/mongo"

type AppData struct {
	Env         *Env
	MongoClient *mongo.Client
	Database    *mongo.Database
}
