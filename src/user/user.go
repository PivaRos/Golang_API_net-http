package user

import (
	"errors"
	"go-api/src/enums"
	"strings"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string             `bson:"fName" json:"fName"  validate:"required"`
	LastName     string             `bson:"lName" json:"lName"  validate:"required"`
	Phone        string             `bson:"phone" json:"phone"  validate:"required"`
	GovId        string             `bson:"govId" json:"govId"  validate:"required"`
	Password     string             `bson:"password" json:"password"`
	Role         enums.Role         `bson:"role" json:"role"  validate:"required"`
	AccessToken  string             `bson:"accessToken" json:"accessToken"  validate:"required"`
	RefreshToken string             `bson:"refreshToken" json:"refreshToken"  validate:"required"`
}

func (c *User) Validate() error {

	validate := validator.New()

	// Validate the struct
	err := validate.Struct(c)
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
