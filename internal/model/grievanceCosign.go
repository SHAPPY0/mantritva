package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrievanceCosign struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	GrievanceID primitive.ObjectID `json:"grievance_id" bson:"grievance_id"`
	UserID      int                `json:"user_id" bson:"user_id"`
	CosignedAt  time.Time          `json:"cosigned_at" bson:"cosigned_at"`
	AuditFieldsBSON `bson:",inline"`
}
