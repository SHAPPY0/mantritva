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
	GetReportKSA() ([]model.KSA, error)
	Create(*model.MemberReport) error
	GetMemberScore(memberID uint, year string) (*model.MemberReport, error)
	GetQuestionnairesByCategory(string) (*[]model.Questionnaire, error)
	GetReactionsByPostID(primitive.ObjectID) ([]model.PostReaction, error)
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

func (r *reportRepository) GetReportKSA() ([]model.KSA, error) {
	var reportKSA []model.KSA
	if err := r.DB.Order("ksa_name ASC").Find(&reportKSA).Error; err != nil {
		return nil, err
	}
	return reportKSA, nil
}

func (r *reportRepository) Create(mr *model.MemberReport) error {
	_, err := r.MongoDB.Insert(r.MongoDB.MemberReportsCollection, mr)
	if err != nil {
		return err
	}
	return nil
}

func (r *reportRepository) GetMemberScore(memberID uint, year string) (*model.MemberReport, error) {
	var result model.MemberReport
	filter := bson.M{
		"member_id": memberID,
		"year":      year,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	rowFound, err := r.MongoDB.FindOne(r.MongoDB.MemberReportsCollection, filter, opts, &result)
	if err != nil {
		return nil, err
	}
	if !rowFound {
		log.Warnf("No report found for memberID %d for year %s", memberID, year)
		result = model.MemberReport{
			MemberID:     uint(memberID),
			Year:         year,
			Developments: model.DevelopmentReport{Weight: 40},
			Grievances:   model.GrievanceReport{Weight: 30},
			NewsPresses:  model.NewsPressReport{Weight: 15},
			FIRs:         model.FIRReport{Weight: 15},
		}
	}
	return &result, nil
}

func (r *reportRepository) GetQuestionnairesByCategory(category string) (*[]model.Questionnaire, error) {
	var questionnaire []model.Questionnaire
	n, ok := r.QuestionnaireCache[category]
	if ok {
		questionnaire = *n
	} else {
		if err := r.DB.Where("category = ? AND is_deleted = false", category).Find(&questionnaire).Error; err != nil {
			return nil, err
		}
		r.QuestionnaireCache[category] = &questionnaire
	}
	return &questionnaire, nil
}

func (r *reportRepository) GetReactionsByPostID(postId primitive.ObjectID) ([]model.PostReaction, error) {
	result := make([]model.PostReaction, 0)
	filter := bson.M{
		"post_id":    postId,
		"is_deleted": false,
	}
	opts := options.Find()
	if err := r.MongoDB.FindMany(r.MongoDB.PostReactionCollection, filter, opts, &result); err != nil {
		return nil, err
	}
	return result, nil
}
