package repository

import (
	"fmt"
	"log"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type NewsPressPostRepository interface {
	Create(*model.NewsPressPost) error
	List() ([]model.NewsPressPost, error)
	GetPostByID(primitive.ObjectID) (*model.NewsPressPost, error)
	UpdatePostVoteCountsByID(primitive.ObjectID, uint, uint) error
	GetReactionByUserIDAndPostID(int, primitive.ObjectID) (*model.PostReaction, error)
	SaveReaction(*model.PostReaction, bool) error
	DeleteReaction(primitive.ObjectID, int, primitive.ObjectID) error
	// UpdateVote(postID string, upvoteChange int, downvoteChange int) error
	ConsoleGetNewsPressList() ([]model.NewsPressPost, error)
	GetPostsByPartyID(pagination dto.Pagination, partyID int) (*[]model.NewsPressPost, error)
}

type newsPressPostRepository struct {
	DB *config.MongoDB
}

func NewNewsPressPostRepository(db *gorm.DB, mongoDB *config.MongoDB) NewsPressPostRepository {
	return &newsPressPostRepository{
		DB: mongoDB,
	}
}

func (r *newsPressPostRepository) Create(np *model.NewsPressPost) error {
	_, err := r.DB.Insert(r.DB.NewsPressPostCollection, np)
	if err != nil {
		return err
	}
	return nil
}

func (r *newsPressPostRepository) List() ([]model.NewsPressPost, error) {
	var newsPressPosts []model.NewsPressPost
	opts := options.Find()
	if err := r.DB.FindMany(r.DB.NewsPressPostCollection, bson.M{}, opts, &newsPressPosts); err != nil {
		return newsPressPosts, err
	}
	return newsPressPosts, nil
}

func (r *newsPressPostRepository) GetPostByID(postID primitive.ObjectID) (*model.NewsPressPost, error) {
	var newsPressPost model.NewsPressPost
	filter := bson.M{"_id": postID, "status": "Published"}
	rowFound, err := r.DB.FindOne(r.DB.NewsPressPostCollection, filter, &newsPressPost)
	if err != nil {
		log.Printf("Error fetching newpresspost by postID %s", postID)
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no newspress post found for postID %s", postID)
	}
	return &newsPressPost, nil
}

func (r *newsPressPostRepository) UpdatePostVoteCountsByID(postID primitive.ObjectID, upvoteCount, donwvoteCount uint) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$set": bson.M{"upvote_count": upvoteCount, "downvote_count": donwvoteCount}}

	_, err := r.DB.UpdateOne(r.DB.NewsPressPostCollection, filter, update)
	if err != nil {
		log.Printf("Error updating vote counts in newsPressPost: %v", err)
		return err
	}
	return nil
}

func (r *newsPressPostRepository) GetReactionByUserIDAndPostID(userID int, postID primitive.ObjectID) (*model.PostReaction, error) {
	var reaction model.PostReaction
	filter := bson.M{"user_id": userID, "post_id": postID, "is_deleted": false}
	rowFound, err := r.DB.FindOne(r.DB.PostReactionCollection, filter, &reaction)
	if err != nil {
		log.Printf("Error fetching reaction: %v", err)
		return nil, err
	}
	if !rowFound {
		return nil, nil
	}
	return &reaction, nil
}

func (r *newsPressPostRepository) DeleteReaction(id primitive.ObjectID, userID int, postID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "post_id": postID}
	update := bson.M{"$set": bson.M{"is_deleted": true, "deleted_at": utils.CreateCurrentDateTime()}}
	_, err := r.DB.UpdateOne(r.DB.PostReactionCollection, filter, update)
	if err != nil {
		log.Printf("Error deleting reaction: %v", err)
		return err
	}

	return nil
}

func (r *newsPressPostRepository) SaveReaction(reaction *model.PostReaction, isNewReaction bool) error {
	_, err := r.DB.Insert(r.DB.PostReactionCollection, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *newsPressPostRepository) ConsoleGetNewsPressList() ([]model.NewsPressPost, error) {
	var newsPressPosts []model.NewsPressPost
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if err := r.DB.FindMany(r.DB.NewsPressPostCollection, bson.M{}, opts, &newsPressPosts); err != nil {
		return newsPressPosts, err
	}
	return newsPressPosts, nil
}

func (r *newsPressPostRepository) GetPostsByPartyID(pagination dto.Pagination, partyID int) (*[]model.NewsPressPost, error) {
	var newsPressPosts []model.NewsPressPost
	filter := bson.M{"party_id": partyID, "status": "Published"}
	opts := options.Find().
		SetLimit(int64(pagination.Limit)).
		SetSkip(int64(pagination.Offset))
	if err := r.DB.FindMany(r.DB.NewsPressPostCollection, filter, opts, &newsPressPosts); err != nil {
		return nil, err
	}
	return &newsPressPosts, nil
}
