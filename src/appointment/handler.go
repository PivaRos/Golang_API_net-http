package appointment

import (
	"encoding/json"
	"go-api/src/enums"
	"go-api/src/utils"
	"net/http"

	"github.com/go-redis/redis/v8"
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

// get all records
func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.RoleContextKey).(enums.Role)
	if !ok {
		http.Error(w, "Role not found", http.StatusInternalServerError)
		return
	}
	userId, ok := r.Context().Value(utils.UserIdContextKey).(string)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	switch role {
	case enums.Admin:
		if id != "" {
			//get single
			appointment, err := h.s.GetById(userId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					http.Error(w, "No document found with the given ID", http.StatusNotFound)
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

		} else {
			//get all of the appointments
			appointment, err := h.s.(userId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					http.Error(w, "No document found with the given ID", http.StatusNotFound)
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
		}
		// Handle Admin role: get all the appointments overall
	case enums.Customer:
		// Handle Customer role: get all its appointments
	case enums.Worker:
		// Handle Worker role: get all its appointments
	default:
		http.Error(w, "Invalid role", http.StatusForbidden)
	}

}

// get single document
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.RoleContextKey).(enums.Role)
	if !ok {
		http.Error(w, "Role not found", http.StatusInternalServerError)
		return
	}
	userId, ok := r.Context().Value(utils.UserIdContextKey).(string)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}

	switch role {
	case enums.Admin:
		customerId := r.URL.Query().Get("customerId")
		workerId := r.URL.Query().Get("workerId")
		if customerId != "" {
			documentId := r.URL.Query().Get("documentId")
			if documentId != "" {
				//get single document of the customer
				appointment, err := h.s.GetByDocumentIdAndCustomerId(documentId, customerId)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						http.Error(w, "No document found", http.StatusNotFound)
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
			} else {
				//get all documents of the customer
				appointments, err := h.s.GetAllByCustomerId(customerId)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						http.Error(w, "No documents found with the given customerId ", http.StatusNotFound)
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
		} else if workerId != "" {
			documentId := r.URL.Query().Get("documentId")
			if documentId != "" {
				//get single document of the worker
				appointment, err := h.s.GetByDocumentIdAndCustomerId(documentId, workerId)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						http.Error(w, "No document found", http.StatusNotFound)
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
			} else {
				//get all documents of the worker
				appointments, err := h.s.GetAllByWorkerId(workerId)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						http.Error(w, "No documents found with the given workerId ", http.StatusNotFound)
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
		} else {
			documentId := r.URL.Query().Get("documentId")
			if documentId != "" {
				//get single document of the worker
				appointment, err := h.s.GetByDocumentId(documentId)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						http.Error(w, "No document found", http.StatusNotFound)
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
			} else {
				//get all the documents
				appointments, err := h.s.GetAll()
				if err != nil {
					if err == mongo.ErrNoDocuments {
						http.Error(w, "No documents found ", http.StatusNotFound)
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
		}
	case enums.Customer:
		// Handle Customer role

	case enums.Worker:
		// Handle Worker role
		customerId := r.URL.Query().Get("customerId")
		if customerId != "" {
			documentId := r.URL.Query().Get("documentId")
			if documentId != "" {
				//get single document of the customer
			} else {
				//get all documents of the customer
			}
		} else {
			documentId := r.URL.Query().Get("documentId")
			if documentId != "" {
				//get single document of the worker
			} else {
				//get all the documents of the current worker
			}
		}
	default:
		http.Error(w, "Invalid role", http.StatusForbidden)
	}

}
