package worker

import (
	"go-api/src/enums"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkHour struct {
	DayOfWeek enums.DayOfWeek `bson:"dayOfWeek" json:"dayOfWeek"`
	StartTime time.Time       `bson:"startTime" json:"startTime"`
	EndTime   time.Time       `bson:"endTime" json:"endTime"`
}

type Worker struct {
	primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	WorkHours          []WorkHour           `bson:"workHours" json:"workHours"`
	AcceptedCares      []primitive.ObjectID `bson:"acceptedCares" json:"acceptedCares"`
}
