package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	Create(*model.Event) error
	List() ([]*model.Event, error)
}

type eventRepository struct {
	DB *config.MongoDB
}

func NewEventRepository(mongoDB *config.MongoDB) EventRepository {
	return &eventRepository{
		DB: mongoDB,
	}
}

func (r *eventRepository) Create(e *model.Event) error {
	_, err := r.DB.Insert(r.DB.EventCollection, e)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) List() ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if err := r.DB.FindMany(r.DB.EventCollection, bson.M{}, opts, &events); err != nil {
		return events, err
	}
	return events, nil
}
