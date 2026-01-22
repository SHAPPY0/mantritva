package repository

import (
	"fmt"
	"log"
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

type NewsPressPostRepository interface {
	Create(*model.NewsPressPost) error
	List() ([]model.NewsPressPost, error)
	GetPostByID(primitive.ObjectID) (*model.NewsPressPost, error)
	ConsoleGetNewsPressPostByID(primitive.ObjectID) (*model.NewsPressPost, error)
	UpdatePostVoteCountsByID(primitive.ObjectID, uint, uint) error
	GetReactionByUserIDAndPostID(int, primitive.ObjectID) (*model.PostReaction, error)
	SaveReaction(*model.PostReaction, bool) error
	DeleteReaction(primitive.ObjectID, int, primitive.ObjectID) error
	// UpdateVote(postID string, upvoteChange int, downvoteChange int) error
	ConsoleGetNewsPressList() ([]model.NewsPressPost, error)
	GetPostsByPartyID(pagination dto.Pagination, partyID int) (*[]model.NewsPressPost, error)
	GetPostsByMemberID(pagination dto.Pagination, memberID int) (*[]model.NewsPressPost, error)
	GetPostsByEntityID(pagination dto.Pagination, entityID int) (*[]model.NewsPressPost, error)
	UpdatePost(string, map[string]interface{}) (bool, error)
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
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.NewsPressPostCollection, filter, opts, &newsPressPost)
	if err != nil {
		log.Printf("Error fetching newpresspost by postID %s", postID)
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no newspress post found for postID %s", postID)
	}
	return &newsPressPost, nil
}

func (r *newsPressPostRepository) ConsoleGetNewsPressPostByID(id primitive.ObjectID) (*model.NewsPressPost, error) {
	var newsPressPost model.NewsPressPost
	filter := bson.M{"_id": id}
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.NewsPressPostCollection, filter, opts, &newsPressPost)
	if err != nil {
		log.Printf("Error fetching newpresspost by postID %s", id)
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("no newspress post found for postID %s", id)
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
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.PostReactionCollection, filter, opts, &reaction)
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
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(pagination.Limit)).
		SetSkip(int64(pagination.Offset))
	if err := r.DB.FindMany(r.DB.NewsPressPostCollection, filter, opts, &newsPressPosts); err != nil {
		return nil, err
	}
	return &newsPressPosts, nil
}

func (r *newsPressPostRepository) GetPostsByMemberID(pagination dto.Pagination, memberID int) (*[]model.NewsPressPost, error) {
	var newsPressPosts []model.NewsPressPost
	filter := bson.M{"member_id": memberID, "status": "Published"}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(pagination.Limit)).
		SetSkip(int64(pagination.Offset))
	if err := r.DB.FindMany(r.DB.NewsPressPostCollection, filter, opts, &newsPressPosts); err != nil {
		return nil, err
	}
	return &newsPressPosts, nil
}

func (r *newsPressPostRepository) GetPostsByEntityID(pagination dto.Pagination, entityID int) (*[]model.NewsPressPost, error) {
	var newsPressPosts []model.NewsPressPost
	filter := bson.M{"entity_id": entityID, "status": "Published"}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(pagination.Limit)).
		SetSkip(int64(pagination.Offset))
	if err := r.DB.FindMany(r.DB.NewsPressPostCollection, filter, opts, &newsPressPosts); err != nil {
		return nil, err
	}
	return &newsPressPosts, nil
}

func (r *newsPressPostRepository) UpdatePost(id string, data map[string]interface{}) (bool, error) {
	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": postID}
	updateFields := bson.M{}

	allowed := map[string]bool{
		"event": true, "subject": true, "type": true, "status": true,
		"tone": true, "statement_date": true, "party_id": true,
		"member_id": true, "entity_id": true, "statement": true, "updated_by": true,
		"updated_at": true,
	}
	for k, v := range data {
		if !allowed[k] {
			continue
		}
		switch k {
		case "statement_date":
			if s, ok := v.(string); ok && s != "" {
				if t, err := time.Parse("2006-01-02", s); err == nil {
					updateFields[k] = t
				} else {
					if t2, err2 := time.Parse(time.RFC3339, s); err2 == nil {
						updateFields[k] = t2
					} else {
						return false, err
					}
				}
			}
		default:
			updateFields[k] = v
		}
	}

	if len(updateFields) == 0 {
		return false, fmt.Errorf("No fields to update")
	}

	update := bson.M{"$set": updateFields}
	_, err = r.DB.UpdateOne(r.DB.NewsPressPostCollection, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}
