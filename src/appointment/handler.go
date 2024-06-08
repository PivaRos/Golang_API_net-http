package appointment

import (
	"encoding/json"
	"go-api/src/care"
	"go-api/src/enums"
	"go-api/src/user"
	"go-api/src/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserDataContextKey).(user.User)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}
	appointments, err := h.s.GetAllByCustomerId(user.Id.Hex())
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	bytes, err := json.Marshal(appointments)
	if err != nil {
		utils.HandleError(w, err)
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
			utils.HandleError(w, err)
			return
		}
		bytes, err := json.Marshal(appointment)
		if err != nil {
			utils.HandleError(w, err)
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
		utils.HandleError(w, err)
		return
	}
	//check if the date did not pass already
	if date.Before(time.Now()) {
		http.Error(w, "date is passed", http.StatusInternalServerError)
		return
	}
	availableHours, err := h.s.GetAvailableTimes(careId, date)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	bytes, err := json.Marshal(availableHours)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {

	var appointment StampAppointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	err = appointment.Validate()
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	user, ok := r.Context().Value(utils.UserDataContextKey).(user.User)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	// validate that the times are actually available
	availableTimes, err := h.s.GetAvailableTimes(appointment.CareId.Hex(), appointment.StartTime)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	careServices := care.CreateServices(h.s.db)
	Care, err := careServices.GetById(appointment.CareId.Hex())
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	var selectedWorkerId string
	*availableTimes = utils.AdjustDates(availableTimes, appointment.StartTime)
	array := *utils.CheckIfTimeCanBeMounted(availableTimes, appointment.StartTime, Care.Duration)
	if len(array) > 0 {
		selectedWorkerId = array[0]
	}
	if selectedWorkerId == "" {
		http.Error(w, "there is no available worker for this job", http.StatusNotFound)
		return
	}
	selectedWorkerObjectId, err := primitive.ObjectIDFromHex(selectedWorkerId)
	if err != nil {
		utils.HandleError(w, err)
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
