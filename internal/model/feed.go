package model

import (
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
