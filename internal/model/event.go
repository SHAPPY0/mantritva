package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Description     string             `json:"description" bson:"description"`
	Date            time.Time          `json:"date" bson:"date"`
	AuditFieldsBSON `bson:",inline"`
}
