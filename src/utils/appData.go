package utils

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppData struct {
	Env         *Env
	MongoClient *mongo.Client
	Database    *mongo.Database
	RedisClient *redis.Client
}
