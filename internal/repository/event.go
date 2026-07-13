package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	Create(*model.Event) error
	List() ([]*model.Event, error)
	GetByID(string) (*model.Event, error)
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

func (r *eventRepository) GetByID(id string) (*model.Event, error) {
	var event model.Event
	opts := options.FindOne()
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Invalid ObjectID: %s", id)
		return nil, err
	}
	ok, err := r.DB.FindOne(r.DB.EventCollection, bson.M{"_id": objId}, opts, &event)
	if err != nil {
		return nil, err
	}
	log.Infof("GetByID: Found event with ID %s: %v", id, ok)
	if !ok {
		return nil, nil
	}
	return &event, nil
}
