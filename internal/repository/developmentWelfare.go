package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type DevelopmentWelfareRepository interface {
	Create(*model.DevelopmentWelfare) error
	List() ([]model.DevelopmentWelfare, error)
	GetPostByID(primitive.ObjectID) (*model.DevelopmentWelfare, error)
	UpdatePostVoteCountsByID(primitive.ObjectID, uint, uint) error
	GetByMemberID(int) ([]model.DevelopmentWelfare, error)
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

func (r *developmentWelfareRepository) GetPostByID(id primitive.ObjectID) (*model.DevelopmentWelfare, error) {
	var devWel model.DevelopmentWelfare
	filter := bson.M{"_id": id}
	opts := options.FindOne()
	if _, err := r.DB.FindOne(r.DB.DevelopmentWelfare, filter, opts, &devWel); err != nil {
		return nil, err
	}
	return &devWel, nil
}

func (r *developmentWelfareRepository) UpdatePostVoteCountsByID(postID primitive.ObjectID, upvoteCount, donwvoteCount uint) error {
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
