package admin

import "net/http"

type Handler struct {
}

func (h *Handler) CreateCare(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
