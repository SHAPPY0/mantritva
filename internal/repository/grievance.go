package repository

import (
	"fmt"
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type GrievanceRepository interface {
	Create(*model.Grievance) error
	ListByUser(string) ([]model.Grievance, error)
	GetByID(primitive.ObjectID) (*model.Grievance, error)
	GetByGrievanceID(string) (*model.Grievance, error)
	ConsoleGetGrievancesList() ([]model.Grievance, error)
	ConsoleCountGrievances() (int64, error)
	ConsoleGetGrievanceByID(string) (*model.Grievance, error)
	ConsoleUpdateGrievanceStatus(string, dto.GrievanceStatus) error
	GetByMemberID(dto.Pagination, int, string) (*[]model.Grievance, error)
	GetByMemberIDAndYear(dto.Pagination, int, int) (*[]model.Grievance, error)
	GetByStateID(int) (*[]model.Grievance, error)
	UpdateVoteCountsByID(primitive.ObjectID, uint, uint) error
	IncrementCosignCountByID(primitive.ObjectID) error
}

type grievanceRepository struct {
	DB *config.MongoDB
}

func NewGrievanceRepository(db *gorm.DB, mongoDB *config.MongoDB) GrievanceRepository {
	return &grievanceRepository{
		DB: mongoDB,
	}
}

func (r *grievanceRepository) Create(g *model.Grievance) error {
	_, err := r.DB.Insert(r.DB.GrievanceCollection, g)
	if err != nil {
		return err
	}
	return nil
}

func (r *grievanceRepository) ListByUser(userID string) ([]model.Grievance, error) {
	filter := map[string]interface{}{
		"created_by": userID,
	}
	opts := options.Find().SetSort(map[string]interface{}{
		"created_at": -1,
	})
	var grievances []model.Grievance
	err := r.DB.FindMany(r.DB.GrievanceCollection, filter, opts, &grievances)
	if err != nil {
		return nil, err
	}
	return grievances, nil
}

func (r *grievanceRepository) GetByID(id primitive.ObjectID) (*model.Grievance, error) {
	filter := map[string]interface{}{
		"_id": id,
	}
	var grievance model.Grievance
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.GrievanceCollection, filter, opts, &grievance)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no grievance found for id %s", id)
	}
	return &grievance, nil
}

func (r *grievanceRepository) GetByGrievanceID(grievanceID string) (*model.Grievance, error) {
	filter := map[string]interface{}{
		"grievance_id": grievanceID,
	}
	var grievance model.Grievance
	opts := options.FindOne()
	_, err := r.DB.FindOne(r.DB.GrievanceCollection, filter, opts, &grievance)
	if err != nil {
		return nil, err
	}
	return &grievance, nil
}

func (r *grievanceRepository) ConsoleGetGrievancesList() ([]model.Grievance, error) {
	var grievances []model.Grievance
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if err := r.DB.FindMany(r.DB.GrievanceCollection, bson.M{}, opts, &grievances); err != nil {
		return grievances, err
	}
	return grievances, nil
}

func (r *grievanceRepository) ConsoleCountGrievances() (int64, error) {
	return r.DB.Count(r.DB.GrievanceCollection, bson.M{})
}

func (r *grievanceRepository) ConsoleGetGrievanceByID(grievanceID string) (*model.Grievance, error) {
	filter := map[string]interface{}{
		"grievance_id": grievanceID,
	}
	var grievance model.Grievance
	opts := options.FindOne()
	_, err := r.DB.FindOne(r.DB.GrievanceCollection, filter, opts, &grievance)
	if err != nil {
		return nil, err
	}
	return &grievance, nil
}

func (r *grievanceRepository) ConsoleUpdateGrievanceStatus(grievanceID string, status dto.GrievanceStatus) error {
	filter := map[string]interface{}{
		"grievance_id": grievanceID,
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
			"final_status": status.Status,
		},
	}
	_, err := r.DB.UpdateOne(r.DB.GrievanceCollection, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *grievanceRepository) GetByMemberID(pagination dto.Pagination, memberID int, position string) (*[]model.Grievance, error) {
	skip := (pagination.Offset - 1) * pagination.Offset
	limit := pagination.Limit
	filter := map[string]interface{}{
		"final_status": bson.M{
			"$ne": utils.UNDERREVIEW,
		},
		"is_deleted": false,
	}
	switch position {
	case "Cabinet Minister":
		filter["cabinet_minister_id"] = memberID
	case "State Minister":
		filter["state_leader_id"] = memberID
	case "Member Of Parliament":
		filter["member_of_parliament_id"] = memberID
	default:
		filter["$or"] = []bson.M{
			{"cabinet_minister_id": memberID},
			{"state_leader_id": memberID},
			{"member_of_parliament_id": memberID},
		}
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	var grievances []model.Grievance
	err := r.DB.FindMany(r.DB.GrievanceCollection, filter, opts, &grievances)
	if err != nil {
		return nil, err
	}
	return &grievances, nil
}

func (r *grievanceRepository) GetByStateID(stateID int) (*[]model.Grievance, error) {
	filter := map[string]interface{}{
		"state_id": stateID,
		"final_status": bson.M{
			"$ne": utils.UNDERREVIEW,
		},
		"is_deleted": false,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	grievances := make([]model.Grievance, 0)
	err := r.DB.FindMany(r.DB.GrievanceCollection, filter, opts, &grievances)
	if err != nil {
		return nil, err
	}
	return &grievances, nil
}

func (r *grievanceRepository) UpdateVoteCountsByID(postID primitive.ObjectID, upvoteCount, donwvoteCount uint) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$set": bson.M{"upvote_count": upvoteCount, "downvote_count": donwvoteCount}}
	_, err := r.DB.UpdateOne(r.DB.GrievanceCollection, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *grievanceRepository) IncrementCosignCountByID(postID primitive.ObjectID) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$inc": bson.M{"cosign_count": 1}}
	_, err := r.DB.UpdateOne(r.DB.GrievanceCollection, filter, update)
	return err
}

func prepareYearFilter(year int) bson.M {
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	return bson.M{"$gte": start, "$lt": end}
}

func (r *grievanceRepository) GetByMemberIDAndYear(pagination dto.Pagination, memberID, year int) (*[]model.Grievance, error) {
	var grievances []model.Grievance
	skip := (pagination.Offset - 1) * pagination.Offset
	limit := pagination.Limit
	filter := bson.M{
		"final_status": bson.M{
			"$ne": utils.UNDERREVIEW,
		},
		"$or": []bson.M{
			{"cabinet_minister_id": memberID},
			{"state_leader_id": memberID},
			{"member_of_parliament_id": memberID},
		},
		"created_at": prepareYearFilter(year),
		"is_deleted": false,
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))
	err := r.DB.FindMany(r.DB.GrievanceCollection, filter, opts, &grievances)
	if err != nil {
		return nil, err
	}
	return &grievances, nil
}
