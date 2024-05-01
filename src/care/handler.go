package care

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

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

type handler struct {
	s *services
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	log.Println(id)
	err := h.s.Delete(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No document found with the given ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete the document: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
