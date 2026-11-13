package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuggestionSupport struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	SuggestionID primitive.ObjectID `json:"suggestion_id" bson:"suggestion_id"`
	UserID       int                `json:"user_id" bson:"user_id"`
	SupportedAt  time.Time          `json:"supported_at" bson:"supported_at"`
	AuditFieldsBSON `bson:",inline"`
}
