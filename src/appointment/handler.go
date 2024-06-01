package appointment

import (
	"encoding/json"
	"go-api/src/user"
	"go-api/src/utils"
	"go-api/src/worker"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	s *services
}

func CreateHandler(db *mongo.Database, redis redis.Client) *handler {
	return &handler{
		s: CreateServices(db, redis),
	}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserDataContextKey).(user.User)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}
	appointments, err := h.s.GetAllByCustomerId(user.Id.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No documents found with the given id", http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	bytes, err := json.Marshal(appointments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserDataContextKey).(user.User)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	if id != "" {

		appointment, err := h.s.GetByDocumentIdAndCustomerId(id, user.Id.Hex())
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No document found with the given id", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		bytes, err := json.Marshal(appointment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
		return
	} else {
		http.Error(w, "Invalid appointment id", http.StatusBadRequest)
		return
	}
}

func (h *handler) GetAvailableTime(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	careId := query.Get("careid")

	dateString := query.Get("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//check if the date did not pass already
	if date.Before(time.Now()) {
		http.Error(w, "date is passed", http.StatusInternalServerError)
		return
	}
	workerCollection := h.s.db.Collection("workers")
	careObjectId, err := primitive.ObjectIDFromHex(careId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cursor, err := workerCollection.Find(r.Context(), bson.M{"certifiedCares": careObjectId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No workers is certified to do this care", http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var workers []worker.Worker
	err = cursor.All(r.Context(), &workers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var objectIds []primitive.ObjectID
	for _, worker := range workers {
		objectIds = append(objectIds, worker.ObjectID)
	}
	appointments, err := h.s.GetAppointmentsByWorkersAndDate(date, objectIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	bytes, err := json.Marshal(availableHoursMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
