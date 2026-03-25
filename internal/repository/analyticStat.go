package repository

import (
	"fmt"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnalyticStatRepository interface {
	Create(*model.AnalyticStat) (any, error)
	Update(*model.AnalyticStat) (bool, error)
	GetByYear(string, int) (*model.AnalyticStat, error)
}

type analyticStatRepository struct {
	DB *config.MongoDB
}

func NewAnalyticStat(mongoDB *config.MongoDB) AnalyticStatRepository {
	return &analyticStatRepository{
		DB: mongoDB,
	}
}

func (r *analyticStatRepository) Create(data *model.AnalyticStat) (any, error) {
	result, err := r.DB.Insert(r.DB.AnalyticStatsCollection, data)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *analyticStatRepository) Update(data *model.AnalyticStat) (bool, error) {
	filter := bson.M{"_id": data.ID, "stat_type": data.StatType, "year": data.Year}
	update := bson.M{"$set": bson.M{"data": data.Data, "metadata.computed_at": time.Now()}}
	_, err := r.DB.UpdateOne(r.DB.AnalyticStatsCollection, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *analyticStatRepository) GetByYear(statType string, year int) (*model.AnalyticStat, error) {
	var stat model.AnalyticStat
	filter := bson.M{"stat_type": statType, "year": year}
	opts := options.FindOne()
	rowFound, err := r.DB.FindOne(r.DB.AnalyticStatsCollection, filter, opts, &stat)
	if err != nil {
		log.Errorf("Error while fetching stat %v", err)
		return nil, err
	}
	if !rowFound {
		return nil, fmt.Errorf("No top scorer stat found %d", year)
	}
	return &stat, nil
}
