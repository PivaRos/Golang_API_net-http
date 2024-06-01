package appointment

import (
	"encoding/json"
	"go-api/src/care"
	"go-api/src/enums"
	"go-api/src/user"
	"go-api/src/utils"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
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
	availableHours, err := h.s.GetAvailableTimes(careId, date)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	bytes, err := json.Marshal(availableHours)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {

	var appointment StampAppointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = appointment.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("here")
	user, ok := r.Context().Value(utils.UserDataContextKey).(user.User)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}
	// validate that the times are actually available
	availableTimes, err := h.s.GetAvailableTimes(appointment.CareId.Hex(), appointment.StartTime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Println("here5")
	careServices := care.CreateServices(h.s.db)
	log.Println(appointment.CareId.Hex())
	Care, err := careServices.GetById(appointment.CareId.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("here6")

	var selectedWorkerId string
	for key, times := range *availableTimes {

		for _, time := range times {
			if time.EndTime.Before(time.StartTime.Add(Care.Duration)) || time.EndTime.Equal(time.StartTime.Add(Care.Duration)) {
				//then this is available time
				selectedWorkerId = key
			}
		}
	}
	selectedWorkerObjectId, err := primitive.ObjectIDFromHex(selectedWorkerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createAppointment := Appointment{
		CareId:     appointment.CareId,
		CustomerId: user.Id,
		WorkerId:   selectedWorkerObjectId,
		StartTime:  appointment.StartTime,
		EndTime:    appointment.StartTime.Add(Care.Duration),
		Status:     enums.Pending,
	}
	err = h.s.Create(createAppointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
