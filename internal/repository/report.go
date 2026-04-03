package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetReportKRA() ([]model.KRA, error)
	Create(*model.Score) error
	CreateScore(*model.Score) error
	CreateScoreSnapshot(*model.ScoreSnapshot) error
	Upsert(*model.Score) error
	GetMemberScore(uint, int) (*model.Score, error)
	GetMemberScoreReports(uint, int, int) ([]*model.Score, error)
	GetMemberOverallScore(uint) ([]*model.Score, error)
	GetQuestionnairesByCategory(string) (*[]model.Questionnaire, error)
	GetScoresByYear(int) ([]*model.Score, error)
}

type reportRepository struct {
	DB                 *gorm.DB
	MongoDB            *config.MongoDB
	QuestionnaireCache map[string]*[]model.Questionnaire
}

func NewReportRepository(db *gorm.DB, mongoDB *config.MongoDB) ReportRepository {
	return &reportRepository{
		DB:                 db,
		MongoDB:            mongoDB,
		QuestionnaireCache: make(map[string]*[]model.Questionnaire),
	}
}

func (r *reportRepository) GetReportKRA() ([]model.KRA, error) {
	var reportKRA []model.KRA
	if err := r.DB.Order("kra_name ASC").Find(&reportKRA).Error; err != nil {
		return nil, err
	}
	return reportKRA, nil
}

func (r *reportRepository) Create(scoreData *model.Score) error {
	_, err := r.MongoDB.Insert(r.MongoDB.ScoreCollection, scoreData)
	if err != nil {
		return err
	}
	return nil
}

func (r *reportRepository) CreateScore(scoreData *model.Score) error {
	memberScore, _ := r.GetMemberScore(scoreData.MemberID, scoreData.Year)
	if memberScore != nil {
		_, err := r.MongoDB.DeleteOne(r.MongoDB.ScoreCollection, bson.M{"member_id": scoreData.MemberID, "year": scoreData.Year})
		if err != nil {
			log.Errorf("CreateScore: Error deleting existing score")
			return err
		}
	}
	scoreData.ID = primitive.NewObjectID()
	ok, err := r.MongoDB.Insert(r.MongoDB.ScoreCollection, scoreData)
	if err != nil {
		return err
	}
	if ok != nil {
		return nil
	}
	return nil
}

func (r *reportRepository) CreateScoreSnapshot(scoreData *model.ScoreSnapshot) error {
	ok, err := r.MongoDB.Insert(r.MongoDB.ScoreSnapshotCollection, scoreData)
	if err != nil {
		return err
	}
	if ok != nil {
		return nil
	}
	return nil
}

func (r *reportRepository) Upsert(scoreData *model.Score) error {
	_, err := r.MongoDB.Upsert(r.MongoDB.ScoreCollection, scoreData)
	if err != nil {
		return err
	}
	return nil
}

func (r *reportRepository) GetMemberScore(memberID uint, year int) (*model.Score, error) {
	result := model.Score{
		MemberID: uint(memberID),
		Year:     year,
		// Developments: model.DevelopmentReport{Weight: 40},
		// Grievances:   model.GrievanceReport{Weight: 30},
		// NewsPresses:  model.NewsPressReport{Weight: 15},
		// CALRecords:   model.CALRecordsReport{Weight: 15},
	}
	filter := bson.M{
		"member_id": memberID,
		"year":      year,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	rowFound, err := r.MongoDB.FindOne(r.MongoDB.ScoreCollection, filter, opts, &result)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		log.Warnf("No report found for memberID %d for year %s", memberID, year)
	}
	return &result, nil
}

func (r *reportRepository) GetMemberOverallScore(memberID uint) ([]*model.Score, error) {
	var scores []*model.Score
	filter := bson.M{
		"member_id": memberID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := r.MongoDB.FindMany(r.MongoDB.ScoreCollection, filter, opts, &scores)
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		log.Warnf("No reports found for memberID %d", memberID)
	}
	return scores, nil
}

func (r *reportRepository) GetQuestionnairesByCategory(feedType string) (*[]model.Questionnaire, error) {
	var questionnaire []model.Questionnaire
	n, ok := r.QuestionnaireCache[feedType]
	if ok {
		questionnaire = *n
	} else {
		if err := r.DB.Where("feed_type = ? AND is_deleted = false", feedType).Find(&questionnaire).Error; err != nil {
			return nil, err
		}
		r.QuestionnaireCache[feedType] = &questionnaire
	}
	return &questionnaire, nil
}

func (r *reportRepository) GetMemberScoreReports(memberID uint, from, to int) ([]*model.Score, error) {
	var scores []*model.Score
	filter := bson.M{
		"member_id":  memberID,
		"year":       bson.M{"$gte": from, "$lte": to},
		"is_deleted": false,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := r.MongoDB.FindMany(r.MongoDB.ScoreCollection, filter, opts, &scores)
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		log.Warnf("No reports found for memberID %d", memberID)
	}
	return scores, nil
}

func (r *reportRepository) GetScoresByYear(year int) ([]*model.Score, error) {
	var scores []*model.Score
	filter := bson.M{
		"year":       year,
		"is_deleted": false,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := r.MongoDB.FindMany(r.MongoDB.ScoreCollection, filter, opts, &scores)
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		log.Warnf("No reports found for year %d", year)
	}

	return scores, nil
}
