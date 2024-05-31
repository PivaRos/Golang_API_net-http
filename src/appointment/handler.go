package appointment

import (
	"go-api/src/user"
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
	user, ok := r.Context().Value(utils.UserDataContextKey).(user.User)
	if !ok {
		http.Error(w, "UserId not found", http.StatusInternalServerError)
		return
	}

}

// get single document
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {

}
