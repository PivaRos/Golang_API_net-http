package appointment

import (
	"errors"
	"go-api/src/enums"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	Id         primitive.ObjectID      `bson:"_id,omitempty" json:"id,omitempty"`
	CareId     primitive.ObjectID      `bson:"careId" json:"careId"  validate:"required"`
	CustomerId primitive.ObjectID      `bson:"customerId" json:"customerId"  validate:"required"`
	WorkerId   primitive.ObjectID      `bson:"workerId" json:"workerId"  validate:"required"`
	StartTime  time.Time               `bson:"startTime" json:"startTime"  validate:"required"`
	EndTime    time.Time               `bson:"endTime" json:"endTime"  validate:"required"`
	Status     enums.AppointmentStatus `bson:"status" json:"status"  validate:"required"`
}

func (a *Appointment) Validate() error {
	validate := validator.New()
	// Validate the struct
	err := validate.Struct(a)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("validation failed")
		}

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field(), err.Tag())
		}
		return errors.New("validation errors: " + strings.Join(validationErrors, ", "))
	}
	return nil
}

type StampAppointment struct {
	StartTime time.Time          `json:"startTime" validate:"required"`
	CareId    primitive.ObjectID `json:"careId" validate:"required"`
}

func (a *StampAppointment) Validate() error {
	validate := validator.New()
	// Validate the struct
	err := validate.Struct(a)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("validation failed")
		}

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field(), err.Tag())
		}
		return errors.New("validation errors: " + strings.Join(validationErrors, ", "))
	}
	return nil
}
