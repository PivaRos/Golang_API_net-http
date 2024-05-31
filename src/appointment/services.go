package appointment

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s services) GetByDocumentId(id string) (*Appointment, error) {
	var appointment Appointment
	AppointmentsDb := s.db.Collection("appointments")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	err = AppointmentsDb.FindOne(context.TODO(), filter).Decode(&appointment)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}
func (s services) GetByDocumentIdAndCustomerId(documentId string, workerId string) (*Appointment, error) {
	var appointment Appointment
	AppointmentsDb := s.db.Collection("appointments")
	documentObjID, err := primitive.ObjectIDFromHex(documentId)
	if err != nil {
		return nil, err
	}
	workerObjID, err := primitive.ObjectIDFromHex(workerId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": documentObjID, "workerId": workerObjID}
	err = AppointmentsDb.FindOne(context.TODO(), filter).Decode(&appointment)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}
func (s services) GetByDocumentIdAndWorkerId(documentId string, customerId string) (*Appointment, error) {
	var appointment Appointment
	AppointmentsDb := s.db.Collection("appointments")
	documentObjID, err := primitive.ObjectIDFromHex(documentId)
	if err != nil {
		return nil, err
	}
	customerObjID, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": documentObjID, "customerId": customerObjID}
	err = AppointmentsDb.FindOne(context.TODO(), filter).Decode(&appointment)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (s *services) GetAllByWorkerId(id string) (*[]Appointment, error) {
	AppointmentsDb := s.db.Collection("appointments")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"workerId": objID}
	cursor, err := AppointmentsDb.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var appointments []Appointment
	if err = cursor.All(context.TODO(), &appointments); err != nil {
		return nil, err
	}
	return &appointments, nil
}
func (s *services) GetAllByCustomerId(id string) (*[]Appointment, error) {
	AppointmentsDb := s.db.Collection("appointments")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"customerId": objID}
	cursor, err := AppointmentsDb.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var appointments []Appointment
	if err = cursor.All(context.TODO(), &appointments); err != nil {
		return nil, err
	}
	return &appointments, nil
}
func (s *services) GetAll() (*[]Appointment, error) {
	AppointmentsDb := s.db.Collection("appointments")
	cursor, err := AppointmentsDb.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var appointments []Appointment
	if err = cursor.All(context.TODO(), &appointments); err != nil {
		return nil, err
	}
	return &appointments, nil
}

func (s *services) GetAppointmentsByWorkersAndDate(date time.Time, arr []primitive.ObjectID) ([]Appointment, error) {
	var conditionsArray []bson.M
	for i := 0; i < len(arr); i++ {
		conditionsArray = append(conditionsArray, bson.M{
			"workerId": arr[i],
		})
	}

	appointmentsCollection := s.db.Collection("appointments")

	// Create date range for the entire day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"$and": []bson.M{
			{"$or": conditionsArray},
			{
				"startTime": bson.M{
					"$gte": startOfDay,
					"$lt":  endOfDay,
				},
			},
		},
	}

	cursor, err := appointmentsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var appointments []Appointment
	err = cursor.All(context.TODO(), &appointments)
	if err != nil {
		return nil, err
	}
	return appointments, nil
}
