package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostVoteRepository interface {
	GetPostByID(primitive.ObjectID) (any, error)
	GetReactionByUserIDAndPostID(int, primitive.ObjectID) (*model.Reaction, error)
	DeleteReaction(primitive.ObjectID, int, primitive.ObjectID) error
	SaveReaction(*model.Reaction, bool) error
	UpdatePostVoteCountsByID(primitive.ObjectID, int, int) error
}
