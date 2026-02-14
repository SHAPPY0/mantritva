package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DevelopmentWelfare struct {
	ID                     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Category               string             `json:"category" bson:"category"`
	Type                   string             `json:"type" bson:"type"`
	PartyID                uint               `json:"party_id" bson:"party_id"`
	MemberID               uint               `json:"member_id" bson:"member_id"`
	Name                   string             `json:"name" bson:"name"`
	Description            string             `json:"description" bson:"description"`
	StateID                *uint              `json:"state_id" bson:"state_id"`
	DistrictID             *uint              `json:"district_id" bson:"district_id"`
	VidhanSabhaID          *uint              `json:"vidhansabha_id" bson:"vidhansabha_id"`
	LokSabhaID             *uint              `json:"loksabha_id" bson:"loksabha_id"`
	Location               *string            `json:"location" bson:"location"`
	AnnouncementDate       time.Time          `json:"announcement_date" bson:"announcement_date"`
	ExpectedCompletionDate *time.Time         `json:"expected_completion_date" bson:"expected_completion_date"`
	ActualCompletionDate   *time.Time         `json:"actual_completion_date" bson:"actual_completion_date"`
	Criteria               *string            `json:"criteria" bson:"criteria"`
	BeneficiariesNum       *uint              `json:"beneficiaries_num" bson:"beneficiaries_num"`
	Scope                  string             `json:"scope" bson:"scope"`
	Status                 string             `json:"status" bson:"status"`
	UpvoteCount            *uint              `json:"upvote_count" bson:"upvote_count"`
	DownvoteCount          *uint              `json:"downvote_count" bson:"downvote_count"`
	Documents              []Document         `json:"documents" bson:"documents"`
	StatusUpdates          []StatusUpdates    `json:"status_updates" bson:"status_updates"`
	AuditFieldsBSON        `bson:",inline"`
}
