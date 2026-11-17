package repository

import (
	"context"
	"fmt"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type PollRepository interface {
	Create(*model.Poll) error
	GetByID(primitive.ObjectID) (*model.Poll, error)
	Vote(*model.PollVote) error
	HasVoted(primitive.ObjectID, int) (*model.PollVote, error)
	Results(primitive.ObjectID) (map[int]int64, error)
	CloseExpired() error
}
type pollRepository struct{ DB *config.MongoDB }

func NewPollRepository(db *config.MongoDB) PollRepository { return &pollRepository{DB: db} }
func (r *pollRepository) Create(p *model.Poll) error {
	_, e := r.DB.Insert(r.DB.PollCollection, p)
	return e
}
func (r *pollRepository) GetByID(id primitive.ObjectID) (*model.Poll, error) {
	var p model.Poll
	found, e := r.DB.FindOne(r.DB.PollCollection, bson.M{"_id": id, "is_deleted": false}, options.FindOne(), &p)
	if e != nil {
		return nil, e
	}
	if !found {
		return nil, fmt.Errorf("poll not found")
	}
	return &p, nil
}
func (r *pollRepository) Vote(v *model.PollVote) error {
	_, e := r.DB.Insert(r.DB.PollVoteCollection, v)
	return e
}
func (r *pollRepository) HasVoted(id primitive.ObjectID, user int) (*model.PollVote, error) {
	var v model.PollVote
	found, e := r.DB.FindOne(r.DB.PollVoteCollection, bson.M{"poll_id": id, "user_id": user}, options.FindOne(), &v)
	if e != nil {
		return nil, e
	}
	if !found {
		return nil, nil
	}
	return &v, nil
}
func (r *pollRepository) Results(id primitive.ObjectID) (map[int]int64, error) {
	cur, e := r.DB.PollVoteCollection.Find(context.Background(), bson.M{"poll_id": id})
	if e != nil {
		return nil, e
	}
	defer cur.Close(context.Background())
	out := map[int]int64{}
	for cur.Next(context.Background()) {
		var v model.PollVote
		if e = cur.Decode(&v); e != nil {
			return nil, e
		}
		out[v.OptionID]++
	}
	return out, cur.Err()
}
func (r *pollRepository) CloseExpired() error {
	_, e := r.DB.PollCollection.UpdateMany(context.Background(), bson.M{"status": "ACTIVE", "ends_at": bson.M{"$lte": time.Now()}}, bson.M{"$set": bson.M{"status": "CLOSED", "updated_at": time.Now()}})
	return e
}

