package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsPressPost struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Event           *string            `json:"event" bson:"event"`
	MemberID        uint               `json:"member_id" bson:"member_id"`
	PartyId         uint               `json:"party_id" bson:"party_id"`
	EntityId        uint               `json:"entity_id" bson:"entity_id"`
	Subject         string             `json:"subject" bson:"subject"`
	Type            string             `json:"type" bson:"type"`
	Tone            string             `json:"tone" bson:"tone"`
	Statement       string             `json:"statement" bson:"statement"`
	StatementDate   time.Time          `json:"statement_date" bson:"statement_date"`
	Status          string             `json:"status" bson:"status"`
	Source          *string            `json:"source" bson:"source"`
	UpvoteCount     uint               `json:"upvote_count" bson:"upvote_count"`
	DownvoteCount   uint               `json:"downvote_count" bson:"downvote_count"`
	AuditFieldsBSON `bson:",inline"`
}
