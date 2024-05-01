package care

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Care struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"  validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Duration    time.Duration      `bson:"duration" json:"duration" validate:"required,gt=0"`
}

func CreateHandler(db *mongo.Database) *handler {
	return &handler{
		db: db,
	}
}

type handler struct {
	db *mongo.Database
}

func (h *handler) ValidateCare(c Care) error {
	validate := validator.New()

	// Validate the struct
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("Validation failed")
		}

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field(), err.Tag())
		}
		return errors.New("Validation errors: " + strings.Join(validationErrors, ", "))
	}
	return nil
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var care Care

	err := json.NewDecoder(r.Body).Decode(&care)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.ValidateCare(care)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	CaresDb := h.db.Collection("Cares")

	res, err := CaresDb.InsertOne(context.TODO(), care)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := res.InsertedID.(primitive.ObjectID); ok {
		w.WriteHeader(http.StatusCreated)
	}
}
