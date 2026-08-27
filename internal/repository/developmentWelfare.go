package repository

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type DevelopmentWelfareRepository interface {
	Create(*model.DevelopmentWelfare) error
	List() ([]model.DevelopmentWelfare, error)
	GetByID(primitive.ObjectID) (*model.DevelopmentWelfare, error)
	UpdateVoteCountsByID(primitive.ObjectID, uint, uint) error
	GetByMemberID(int) ([]model.DevelopmentWelfare, error)
	GetByMemberIDAndYear(int, int) ([]model.DevelopmentWelfare, error)
	ConsoleGetDevelopmentWelfaresList() ([]model.DevelopmentWelfare, error)
	ConsoleCountDevelopmentWelfares() (int64, error)
	ConsoleGetDevelopmentWelfareByID(string) (*model.DevelopmentWelfare, error)
	ConsoleUpdateDevelopmentWelfareStatus(string, dto.DevelopmentWelfareStatus) error
}

type developmentWelfareRepository struct {
	DB *config.MongoDB
}

func NewDevelopmentWelfareRepository(db *gorm.DB, mongoDB *config.MongoDB) DevelopmentWelfareRepository {
	return &developmentWelfareRepository{
		DB: mongoDB,
	}
}

func (r *developmentWelfareRepository) Create(d *model.DevelopmentWelfare) error {
	_, err := r.DB.Insert(r.DB.DevelopmentWelfare, d)
	if err != nil {
		return err
	}
	return nil
}

func (r *developmentWelfareRepository) List() ([]model.DevelopmentWelfare, error) {
	var developments []model.DevelopmentWelfare
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.DevelopmentWelfare, bson.M{}, opts, &developments); err != nil {
		return developments, err
	}
	return developments, nil
}

func (r *developmentWelfareRepository) GetByID(id primitive.ObjectID) (*model.DevelopmentWelfare, error) {
	var devWel model.DevelopmentWelfare
	filter := bson.M{"_id": id}
	opts := options.FindOne()
	if _, err := r.DB.FindOne(r.DB.DevelopmentWelfare, filter, opts, &devWel); err != nil {
		log.Errorf("Error while fetching development post %v", err)
		return nil, err
	}
	return &devWel, nil
}

func (r *developmentWelfareRepository) UpdateVoteCountsByID(postID primitive.ObjectID, upvoteCount, donwvoteCount uint) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$set": bson.M{"upvote_count": upvoteCount, "downvote_count": donwvoteCount}}
	_, err := r.DB.UpdateOne(r.DB.DevelopmentWelfare, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *developmentWelfareRepository) GetByMemberID(memberID int) ([]model.DevelopmentWelfare, error) {
	developments := make([]model.DevelopmentWelfare, 0)
	filter := bson.M{"member_id": memberID}
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.DevelopmentWelfare, filter, opts, &developments); err != nil {
		return developments, err
	}
	return developments, nil
}

func (r *developmentWelfareRepository) GetByMemberIDAndYear(memberID, year int) ([]model.DevelopmentWelfare, error) {
	developments := make([]model.DevelopmentWelfare, 0)
	filter := bson.M{
		"member_id":  memberID,
		"created_at": prepareYearFilter(year),
		"is_deleted": false,
	}
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.DevelopmentWelfare, filter, opts, &developments); err != nil {
		return developments, err
	}
	return developments, nil
}

func (r *developmentWelfareRepository) ConsoleGetDevelopmentWelfaresList() ([]model.DevelopmentWelfare, error) {
	var developments []model.DevelopmentWelfare
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if err := r.DB.FindMany(r.DB.DevelopmentWelfare, bson.M{}, opts, &developments); err != nil {
		return developments, err
	}
	return developments, nil
}

func (r *developmentWelfareRepository) ConsoleCountDevelopmentWelfares() (int64, error) {
	return r.DB.Count(r.DB.DevelopmentWelfare, bson.M{})
}

func (r *developmentWelfareRepository) ConsoleGetDevelopmentWelfareByID(id string) (*model.DevelopmentWelfare, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(objectID)
}

func (r *developmentWelfareRepository) ConsoleUpdateDevelopmentWelfareStatus(id string, status dto.DevelopmentWelfareStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
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
	_, err = r.DB.UpdateOne(r.DB.DevelopmentWelfare, filter, update)
	return err
}
