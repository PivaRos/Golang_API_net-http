package appointment

import (
	"context"
	"errors"
	"go-api/src/utils"
	"go-api/src/worker"
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
func (s *services) GetAvailableTimes(careId string, date time.Time) (*map[string][]utils.Times, error) {
	workerCollection := s.db.Collection("workers")
	careObjectId, err := primitive.ObjectIDFromHex(careId)
	if err != nil {
		return nil, err
	}
	cursor, err := workerCollection.Find(context.TODO(), bson.M{"certifiedCares": careObjectId})
	if err != nil {
		return nil, err
	}
	var workers []worker.Worker
	err = cursor.All(context.TODO(), &workers)
	if err != nil {
		return nil, err
	}
	var objectIds []primitive.ObjectID
	for _, worker := range workers {
		objectIds = append(objectIds, worker.ObjectID)
	}
	if len(objectIds) == 0 {
		return nil, errors.New("no workers found")
	}
	appointments, err := s.GetAppointmentsByWorkersAndDate(date, objectIds)
	if err != nil {
		return nil, err
	}
	// get the raw available times

	availableHoursMap := make(map[string][]utils.Times)
	for _, w := range workers {
		for _, hours := range w.WorkHours {
			if hours.WeekDay == date.Weekday() {
				availableHoursMap[w.Hex()] = append(availableHoursMap[w.Hex()], utils.Times{StartTime: hours.StartTime, EndTime: hours.EndTime})
			}
		}
	}

	//update the raw available times
	for _, appointment := range appointments {
		availableHoursMap[string(appointment.WorkerId.Hex())] = utils.SubtractTimes(availableHoursMap[appointment.WorkerId.Hex()], utils.Times{StartTime: appointment.StartTime, EndTime: appointment.EndTime})
	}
	return &availableHoursMap, nil
}

func (s *services) Create(appointment Appointment) error {
	appointmentsCollection := s.db.Collection("appointments")
	res, err := appointmentsCollection.InsertOne(context.TODO(), appointment)
	if err != nil {
		return err
	}
	if _, ok := res.InsertedID.(primitive.ObjectID); ok {
		return nil
	} else {
		return errors.New("unable to complete action")
	}

}

func (s *services) Update(appointmentId primitive.ObjectID, updateData Appointment) error {
	appointmentsCollection := s.db.Collection("appointments")

	// Create the update document
	update := bson.M{
		"$set": bson.M{
			"careId":     updateData.CareId,
			"customerId": updateData.CustomerId,
			"workerId":   updateData.WorkerId,
			"startTime":  updateData.StartTime,
			"endTime":    updateData.EndTime,
			"status":     updateData.Status,
		},
	}

	// Find the appointment by ID and update it
	filter := bson.M{"_id": appointmentId}
	res, err := appointmentsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("appointment not found")
	}

	if res.ModifiedCount == 0 {
		return errors.New("no changes made to the appointment")
	}

	return nil
}
