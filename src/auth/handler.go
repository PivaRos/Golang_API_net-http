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

func (h *handler) SendOTP(w http.ResponseWriter, r *http.Request) {
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
	token, err := h.s.SendOTP(loginCredentials)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		utils.HandleError(w, err)
		return
	}
	tokenRaw, err := json.Marshal(token)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenRaw)
}
func (h *handler) ValidateOTP(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	code := r.URL.Query().Get("otp")
	tokens, err := h.s.ValidateOTP(tokenString, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenRaw, err := json.Marshal(tokens)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenRaw)

}
