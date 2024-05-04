package user

import (
	"go-api/src/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string             `bson:"fName" json:"fName"  validate:"required"`
	LastName     string             `bson:"lName" json:"lName"  validate:"required"`
	Phone        string             `bson:"phone" json:"phone"  validate:"required"`
	GovId        string             `bson:"govId" json:"govId"  validate:"required"`
	Password     string             `bson:"password" json:"password"  validate:"required"`
	Role         enums.Role         `bson:"role" json:"role"  validate:"required"`
	AccessToken  string             `bson:"accessToken" json:"accessToken"  validate:"required"`
	RefreshToken string             `bson:"refreshToken" json:"refreshToken"  validate:"required"`
}
