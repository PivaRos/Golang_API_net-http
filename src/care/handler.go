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
		s: CreateServices(db),
	}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id != "" {
		//get single
		care, err := h.s.GetById(id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No document found with the given ID", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		bytes, err := json.Marshal(care)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)

	} else {
		//get all
		cares, err := h.s.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(cares)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
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

	id := r.PathValue("id")
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
