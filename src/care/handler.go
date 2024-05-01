package care

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Care struct {
	id       string        `bson:"_id" json:"id"`
	date     string        `bson:"date" json:"date"`
	duration time.Duration `bson:"duration" json:"duration"`
}

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
