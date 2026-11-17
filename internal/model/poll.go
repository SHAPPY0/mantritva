package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Poll struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Question        string             `json:"question" bson:"question"`
	PollType        string             `json:"poll_type" bson:"poll_type"`
	Options         []PollOption       `json:"options" bson:"options"`
	DurationDays    int                `json:"duration_days" bson:"duration_days"`
	StartsAt        time.Time          `json:"starts_at" bson:"starts_at"`
	EndsAt          time.Time          `json:"ends_at" bson:"ends_at"`
	Status          string             `json:"status" bson:"status"`
	FeedID          primitive.ObjectID `json:"feed_id" bson:"feed_id"`
	AuditFieldsBSON `bson:",inline"`
}

type PollOption struct {
	ID   int    `json:"id" bson:"id"`
	Text string `json:"text" bson:"text"`
}

type PollVote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PollID    primitive.ObjectID `json:"poll_id" bson:"poll_id"`
	OptionID  int                `json:"option_id" bson:"option_id"`
	UserID    int                `json:"-" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

