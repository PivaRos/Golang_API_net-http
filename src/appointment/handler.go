package appointment

import (
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

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.RoleContextKey).(enums.Role)
	if !ok {
		http.Error(w, "Role not found", http.StatusInternalServerError)
		return
	}
	switch role {
	case enums.Admin:
		// Handle Admin role: get all the appointments overall
	case enums.Customer:
		// Handle Customer role: get all its appointments
	case enums.Worker:
		// Handle Worker role: get all its appointments
	default:
		http.Error(w, "Invalid role", http.StatusForbidden)
	}

}
