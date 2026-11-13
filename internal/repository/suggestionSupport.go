package repository

import (
	"fmt"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SuggestionSupportRepository interface {
	Create(*model.SuggestionSupport) error
	GetByUserAndSuggestionID(int, primitive.ObjectID) (*model.SuggestionSupport, error)
	CountBySuggestionID(primitive.ObjectID) (int64, error)
}

type suggestionSupportRepository struct {
	DB *config.MongoDB
}

func NewSuggestionSupportRepository(mongoDB *config.MongoDB) SuggestionSupportRepository {
	return &suggestionSupportRepository{
		DB: mongoDB,
	}
}

func (r *suggestionSupportRepository) Create(support *model.SuggestionSupport) error {
	_, err := r.DB.Insert(r.DB.SuggestionSupportCollection, support)
	return err
}

func (r *suggestionSupportRepository) GetByUserAndSuggestionID(userID int, suggestionID primitive.ObjectID) (*model.SuggestionSupport, error) {
	filter := bson.M{
		"user_id":       userID,
		"suggestion_id": suggestionID,
		"is_deleted":    false,
	}
	var support model.SuggestionSupport
	rowFound, err := r.DB.FindOne(r.DB.SuggestionSupportCollection, filter, options.FindOne(), &support)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no suggestion support found")
	}
	return &support, nil
}

func (r *suggestionSupportRepository) CountBySuggestionID(suggestionID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"suggestion_id": suggestionID,
		"is_deleted":    false,
	}
	return r.DB.Count(r.DB.SuggestionSupportCollection, filter)
}
