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

func (s *services) GetById(id string) (*Care, error) {
	var care Care
	CaresDb := s.db.Collection("cares")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}

	// Find the document using the filter
	err = CaresDb.FindOne(context.TODO(), filter).Decode(&care)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, errors.New("no care found with this id")
		} else {

			return nil, err
		}
	}

	return &care, nil
}

func (s *services) GetAll() (*[]Care, error) {
	Cares := s.db.Collection("cares")

	cursor, err := Cares.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var cares []Care
	if err = cursor.All(context.TODO(), &cares); err != nil {
		return nil, err
	}
	return &cares, nil
}

func (s *services) Create(care Care) error {
	Cares := s.db.Collection("cares")

	res, err := Cares.InsertOne(context.TODO(), care)
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
	Cares := s.db.Collection("cares")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	result, err := Cares.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
