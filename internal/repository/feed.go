package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedRepository interface {
	Create(*model.Feed) (any, error)
	GetAll() ([]model.Feed, error)
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
