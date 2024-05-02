package auth

import (
	"encoding/json"
	"go-api/src/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateHandler(app *utils.AppData) *handler {
	return &handler{
		s: CreateServices(app),
	}
}

type handler struct {
	s *services
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginCredentials Login
	err := json.NewDecoder(r.Body).Decode(&loginCredentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = loginCredentials.ValidateLogin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := h.s.Login(loginCredentials)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokenRaw, err := json.Marshal(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(tokenRaw)
}
