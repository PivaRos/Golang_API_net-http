package appointment

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateServices(db *mongo.Database, redis redis.Client) *services {
	return &services{
		db:    db,
		redis: redis,
	}
}

type services struct {
	db    *mongo.Database
	redis redis.Client
}


