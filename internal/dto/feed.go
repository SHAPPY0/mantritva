package dto

import (
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedReq struct {
	Type   utils.FeedType     `json:"type"`
	PostID primitive.ObjectID `json:"post_id"`
	Status utils.PostStatus   `json:"status"`
}

type FeedRes struct {
	ID   primitive.ObjectID `json:"id"`
	Type string             `json:"type"`
	Feed any                `json:"feed"`
}

type ReactionReq struct {
	PostID   primitive.ObjectID     `json:"post_id" binding:"required"`
	UserID   int                    `json:"user_id"`
	Reaction utils.PostReactionType `json:"reaction" binding:"required"`
	FeedType utils.FeedType         `json:"feed_type" binding:"required"`
}

type PostVotes struct {
	ID            primitive.ObjectID
	UpvoteCount   int
	DownvoteCount int
}
