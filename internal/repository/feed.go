package repository

import (
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedRepository interface {
	Create(*model.Feed) (any, error)
	UpdateStatus(string, dto.FeedReq) (bool, error)
	GetAll() ([]model.Feed, error)
	GetReactionByUserIDAndPostID(int, primitive.ObjectID) (*model.PostReaction, error)
	SaveReaction(*model.PostReaction, bool) error
	DeleteReaction(primitive.ObjectID, int, primitive.ObjectID) error
	// UpdatePostVoteCountsByID(primitive.ObjectID, int, int) error
}

type feedRepository struct {
	DB *config.MongoDB
}

func NewFeedRepository(db *config.MongoDB) FeedRepository {
	return &feedRepository{
		DB: db,
	}
}

func (r *feedRepository) Create(feed *model.Feed) (any, error) {
	result, err := r.DB.Insert(r.DB.FeedCollection, feed)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *feedRepository) GetAll() ([]model.Feed, error) {
	var feeds []model.Feed
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(50)
	if err := r.DB.FindMany(r.DB.FeedCollection, bson.M{"status": utils.Published}, opts, &feeds); err != nil {
		return nil, err
	}
	return feeds, nil
}

func (r *feedRepository) UpdateStatus(updatedBy string, feed dto.FeedReq) (bool, error) {
	filter := bson.M{"post_id": feed.PostID, "type": feed.Type}
	update := bson.M{"$set": bson.M{"status": feed.Status, "updated_by": updatedBy, "updated_at": time.Now()}}
	_, err := r.DB.UpdateOne(r.DB.FeedCollection, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *feedRepository) GetReactionByUserIDAndPostID(userID int, postID primitive.ObjectID) (*model.PostReaction, error) {
	var reaction model.PostReaction
	filter := bson.M{"user_id": userID, "post_id": postID, "is_deleted": false}
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.PostReactionCollection, filter, opts, &reaction)
	if err != nil {
		log.Errorf("Error fetching reaction: %v", err)
		return nil, err
	}
	if !rowFound {
		return nil, nil
	}
	return &reaction, nil
}

func (r *feedRepository) DeleteReaction(id primitive.ObjectID, userID int, postID primitive.ObjectID) error {
	filter := bson.M{"_id": id, "user_id": userID, "post_id": postID, "is_deleted": false}
	update := bson.M{"$set": bson.M{"is_deleted": true, "deleted_at": utils.CreateCurrentDateTime()}}
	_, err := r.DB.UpdateOne(r.DB.PostReactionCollection, filter, update)
	if err != nil {
		log.Errorf("Error deleting reaction: %v", err)
		return err
	}

	return nil
}

func (r *feedRepository) SaveReaction(reaction *model.PostReaction, isNewReaction bool) error {
	_, err := r.DB.Insert(r.DB.PostReactionCollection, reaction)
	if err != nil {
		return err
	}
	return nil
}
