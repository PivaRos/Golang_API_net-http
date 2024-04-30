package care

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateHandler(db *mongo.Database) *handler {
	return &handler{
		db: db,
	}
}

type handler struct {
	db *mongo.Database
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {

}
