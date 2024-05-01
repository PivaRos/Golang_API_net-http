package care

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
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

func (s *services) Delete(id string) error {
	CaresDb := s.db.Collection("Cares")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	result, err := CaresDb.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
