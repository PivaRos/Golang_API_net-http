package appointment

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	s *services
}

func CreateHandler(db *mongo.Database) *handler {
	return &handler{
		s: CreateServices(db),
	}
}


