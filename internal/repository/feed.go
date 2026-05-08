package repository

import (
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedRepository interface {
	Create(*model.Feed) (any, error)
	UpdateStatus(string, dto.FeedReq) (bool, error)
	GetAll() ([]model.Feed, error)
	GetByID(string) (*model.Feed, error)
	GetByPagination(dto.Pagination) ([]model.Feed, error)
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

func (r *feedRepository) GetByPagination(pagination dto.Pagination) ([]model.Feed, error) {
	var feeds []model.Feed
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(pagination.Offset)).
		SetLimit(int64(pagination.Limit))
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

func (r *feedRepository) GetByID(feedID string) (*model.Feed, error) {
	var feed model.Feed
	id, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id, "status": utils.Published, "is_deleted": false}
	if _, err := r.DB.FindOne(r.DB.FeedCollection, filter, nil, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}
