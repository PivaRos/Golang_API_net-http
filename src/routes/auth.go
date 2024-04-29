package routes

import "net/http"

type Auth struct {
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
