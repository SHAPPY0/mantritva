package repository

import (
	"fmt"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type SuggestionRepository interface {
	Create(*model.Suggestion) error
	ListByUser(string) ([]model.Suggestion, error)
	GetBySuggestionID(string) (*model.Suggestion, error)
	GetByID(primitive.ObjectID) (*model.Suggestion, error)
	IncrementViewCount(primitive.ObjectID) error
	ConsoleGetSuggestionsList() ([]model.Suggestion, error)
	ConsoleGetSuggestionByID(string) (*model.Suggestion, error)
	ConsoleUpdateSuggestionStatus(string, dto.SuggestionStatusUpdate) error
	UpdateVoteCountsByID(primitive.ObjectID, uint, uint) error
}

type suggestionRepository struct {
	DB *config.MongoDB
}

func NewSuggestionRepository(db *gorm.DB, mongoDB *config.MongoDB) SuggestionRepository {
	return &suggestionRepository{
		DB: mongoDB,
	}
}

func (r *suggestionRepository) Create(suggestion *model.Suggestion) error {
	_, err := r.DB.Insert(r.DB.SuggestionCollection, suggestion)
	return err
}

func (r *suggestionRepository) ListByUser(userID string) ([]model.Suggestion, error) {
	filter := bson.M{
		"created_by": userID,
		"is_deleted": false,
	}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	var suggestions []model.Suggestion
	if err := r.DB.FindMany(r.DB.SuggestionCollection, filter, opts, &suggestions); err != nil {
		return nil, err
	}
	return suggestions, nil
}

func (r *suggestionRepository) GetBySuggestionID(suggestionID string) (*model.Suggestion, error) {
	filter := bson.M{
		"suggestion_id": suggestionID,
		"is_deleted":    false,
	}
	var suggestion model.Suggestion
	rowFound, err := r.DB.FindOne(r.DB.SuggestionCollection, filter, options.FindOne(), &suggestion)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no suggestion found for id %s", suggestionID)
	}
	return &suggestion, nil
}

func (r *suggestionRepository) GetByID(id primitive.ObjectID) (*model.Suggestion, error) {
	filter := bson.M{
		"_id":        id,
		"is_deleted": false,
	}
	var suggestion model.Suggestion
	rowFound, err := r.DB.FindOne(r.DB.SuggestionCollection, filter, options.FindOne(), &suggestion)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no suggestion found for id %s", id.Hex())
	}
	return &suggestion, nil
}

func (r *suggestionRepository) IncrementViewCount(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"view_count": 1}}
	_, err := r.DB.UpdateOne(r.DB.SuggestionCollection, filter, update)
	return err
}

func (r *suggestionRepository) ConsoleGetSuggestionsList() ([]model.Suggestion, error) {
	var suggestions []model.Suggestion
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if err := r.DB.FindMany(r.DB.SuggestionCollection, bson.M{}, opts, &suggestions); err != nil {
		return suggestions, err
	}
	return suggestions, nil
}

func (r *suggestionRepository) ConsoleGetSuggestionByID(suggestionID string) (*model.Suggestion, error) {
	filter := bson.M{
		"suggestion_id": suggestionID,
	}
	var suggestion model.Suggestion
	_, err := r.DB.FindOne(r.DB.SuggestionCollection, filter, options.FindOne(), &suggestion)
	if err != nil {
		return nil, err
	}
	return &suggestion, nil
}

func (r *suggestionRepository) ConsoleUpdateSuggestionStatus(
	suggestionID string,
	status dto.SuggestionStatusUpdate,
) error {
	filter := bson.M{
		"suggestion_id": suggestionID,
	}
	update := bson.M{
		"$push": bson.M{
			"status_revisions": bson.M{
				"date":   status.Date.Time,
				"status": status.Status,
				"remark": status.Remark,
				"source": status.Source,
			},
		},
		"$set": bson.M{
			"status": status.Status,
		},
	}
	_, err := r.DB.UpdateOne(r.DB.SuggestionCollection, filter, update)
	return err
}

func (r *suggestionRepository) UpdateVoteCountsByID(postID primitive.ObjectID, upvoteCount, downvoteCount uint) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$set": bson.M{"upvote_count": upvoteCount, "downvote_count": downvoteCount}}
	_, err := r.DB.UpdateOne(r.DB.SuggestionCollection, filter, update)
	return err
}
