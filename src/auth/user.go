package auth

import (
	"go-api/src/role"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccessToken  string             `bson:"accessToken" json:"accessToken"  validate:"required"`
	RefreshToken string             `bson:"refreshToken" json:"refreshToken"  validate:"required"`
	Phone        string             `bson:"phone" json:"phone"  validate:"required"`
	Password     string             `bson:"password" json:"password"  validate:"required"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"  validate:"required"`
	LastModified time.Time          `bson:"lastModified" json:"lastModified"  validate:"required"`
	Role         role.Role          `bson:"role" json:"role"  validate:"required"`
}
