package auth

import "go.mongodb.org/mongo-driver/mongo"

func CreateServices(db *mongo.Database) *services {
	return &services{
		db: db,
	}
}

type services struct {
	db *mongo.Database
}
