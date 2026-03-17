package model

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Type            string             `json:"type" bson:"type"`
	PostID          primitive.ObjectID `json:"post_id" bson:"post_id"`
	Status          utils.PostStatus   `json:"status" bson:"status"`
	AuditFieldsBSON `bson:",inline"`
}

type Reaction struct {
	ID                  primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	PostID              primitive.ObjectID `json:"post_id" bson:"post_id"`
	UserID              int                `json:"user_id" bson:"user_id"`
	Reaction            string             `json:"reaction" bson:"reaction"`
	UserIP              string             `json:"user_ip" bson:"user_ip"`
	SelectedOptions     []int              `json:"selected_options" bson:"selected_options"`
	PrevNormalizedScore *int               `json:"prev_normalized_score" bson:"prev_normalized_score"`
	NormalizedScore     int                `json:"normalized_score" bson:"normalized_score"`
	LockedAt            *time.Time         `json:"locked_at,omitempty" bson:"locked_at,omitempty"`
	Locked              bool               `json:"locked" bson:"locked"`
	AuditFieldsBSON     `bson:",inline"`
}

type ReactionQueue struct {
	PostID          primitive.ObjectID `json:"post_id" bson:"post_id"`
	UserID          int                `json:"user_id" bson:"user_id"`
	Reaction        string             `json:"reaction" bson:"reaction"`
	FeedType        string             `json:"feed_type" bson:"feed_type"`
	SelectedOptions []int              `json:"selected_options" bson:"selected_options"`
	UserIP          string             `json:"user_ip" bson:"user_ip"`
	Processed       bool               `json:"processed" bson:"processed"`
}
