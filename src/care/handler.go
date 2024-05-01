package care

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	s *services
}

func CreateHandler(db *mongo.Database) *handler {
	return &handler{
		s: &services{
			db: db,
		},
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var care Care

	err := json.NewDecoder(r.Body).Decode(&care)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = care.ValidateCare()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.Create(care)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {

}
