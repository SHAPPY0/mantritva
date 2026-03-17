package repository

import (
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReactionRepository interface {
	Create(*model.Reaction, bool) error
	Delete(primitive.ObjectID, int, primitive.ObjectID, bool) error
	GetByPostID(primitive.ObjectID) ([]model.Reaction, error)
	GetByUserIDAndPostID(int, primitive.ObjectID) (*model.Reaction, error)
	GetLockedReactions(primitive.ObjectID) ([]*model.Reaction, error)
	GetAndLockMaturedReactions(time.Time) ([]model.Reaction, error)
}

type reactionRepository struct {
	DB *config.MongoDB
}

func NewReactionRepository(mongoDB *config.MongoDB) ReactionRepository {
	return &reactionRepository{
		DB: mongoDB,
	}
}

func (r *reactionRepository) Create(reaction *model.Reaction, isNewReaction bool) error {
	_, err := r.DB.Insert(r.DB.ReactionCollection, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *reactionRepository) Delete(id primitive.ObjectID, userID int, postID primitive.ObjectID, hardDelete bool) error {
	filter := bson.M{"_id": id, "user_id": userID, "post_id": postID, "is_deleted": false}
	if !hardDelete {
		update := bson.M{"$set": bson.M{"is_deleted": true, "deleted_at": utils.CreateCurrentDateTime()}}
		_, err := r.DB.UpdateOne(r.DB.ReactionCollection, filter, update)
		if err != nil {
			log.Errorf("Error deleting reaction: %v", err)
			return err
		}
	} else {
		_, err := r.DB.DeleteOne(r.DB.ReactionCollection, filter)
		if err != nil {
			log.Errorf("Error hard deleting reaction: %v", err)
			return err
		}
	}

	return nil
}

func (r *reactionRepository) GetByPostID(postId primitive.ObjectID) ([]model.Reaction, error) {
	result := make([]model.Reaction, 0)
	filter := bson.M{"post_id": postId, "is_deleted": false}
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.ReactionCollection, filter, opts, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *reactionRepository) GetByUserIDAndPostID(userID int, postID primitive.ObjectID) (*model.Reaction, error) {
	var reaction model.Reaction
	filter := bson.M{"user_id": userID, "post_id": postID, "is_deleted": false}
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.ReactionCollection, filter, opts, &reaction)
	if err != nil {
		log.Errorf("Error fetching reaction: %v", err)
		return nil, err
	}
	if !rowFound {
		return nil, nil
	}
	return &reaction, nil
}

func (r *reactionRepository) GetLockedReactions(postID primitive.ObjectID) ([]*model.Reaction, error) {
	var reactions []*model.Reaction
	currentDateTime := time.Now()
	hourStart := currentDateTime.Truncate(time.Hour)
	hourEnd := currentDateTime.Truncate(time.Hour).Add(time.Hour).Add(-time.Nanosecond)

	filter := bson.M{
		"post_id":   postID,
		"locked_at": bson.M{"$gte": hourStart.UTC(), "$lte": hourEnd.UTC()},
	}
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.ReactionCollection, filter, opts, &reactions); err != nil {
		log.Errorf("GetLockedReactions: Error fetching reactions: %v", err)
		return nil, err
	}
	return reactions, nil
}

func (r *reactionRepository) GetAndLockMaturedReactions(computeDateTime time.Time) ([]model.Reaction, error) {
	var reactions []model.Reaction
	hourStart := computeDateTime.Truncate(time.Hour)
	hourEnd := computeDateTime.Truncate(time.Hour).Add(time.Hour)
	filter := bson.M{
		"locked_at": bson.M{"$gte": hourStart.UTC(), "$lt": hourEnd.UTC()},
		// "$or":       []bson.M{{"locked": true}, {"locked": nil}},
	}
	log.Infof("Filter: %v", filter)
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.ReactionCollection, filter, opts, &reactions); err != nil {
		log.Errorf("GetAndLockMaturedReactions: Error fetching locked reactions: %v", err)
		return nil, err
	}
	return reactions, nil
}
