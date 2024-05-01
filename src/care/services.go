package care

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateServices(db *mongo.Database) *services {
	return &services{
		db: db,
	}
}

type services struct {
	db *mongo.Database
}

func (s *services) Create(care Care) error {
	CaresDb := s.db.Collection("Cares")

	res, err := CaresDb.InsertOne(context.TODO(), care)
	if err != nil {
		return err
	}
	if _, ok := res.InsertedID.(primitive.ObjectID); ok {
		return nil
	} else {
		return errors.New("unable to complete action")
	}
}
