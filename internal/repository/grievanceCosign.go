package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GrievanceCosignRepository interface {
	Create(*model.GrievanceCosign) error
	GetByUserIDAndGrievanceID(int, primitive.ObjectID) (*model.GrievanceCosign, error)
}

type grievanceCosignRepository struct {
	DB *config.MongoDB
}

func NewGrievanceCosignRepository(mongoDB *config.MongoDB) GrievanceCosignRepository {
	return &grievanceCosignRepository{
		DB: mongoDB,
	}
}

func (r *grievanceCosignRepository) Create(cosign *model.GrievanceCosign) error {
	_, err := r.DB.Insert(r.DB.GrievanceCosignCollection, cosign)
	return err
}

func (r *grievanceCosignRepository) GetByUserIDAndGrievanceID(userID int, grievanceID primitive.ObjectID) (*model.GrievanceCosign, error) {
	var cosign model.GrievanceCosign
	filter := bson.M{
		"user_id":      userID,
		"grievance_id": grievanceID,
		"is_deleted":   false,
	}
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.GrievanceCosignCollection, filter, opts, &cosign)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		return nil, nil
	}
	return &cosign, nil
}
