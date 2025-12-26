package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetReportKSA() ([]model.KSA, error)
	Create(*model.MemberReport) error
	GetMemberScore(memberID uint, year string) (*model.MemberReport, error)
}

type reportRepository struct {
	DB      *gorm.DB
	MongoDB *config.MongoDB
}

func NewReportRepository(db *gorm.DB, mongoDB *config.MongoDB) ReportRepository {
	return &reportRepository{
		DB:      db,
		MongoDB: mongoDB,
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
		log.Infof("No report found for memberID %d for year %s", memberID, year)
		return &result, nil
	}
	log.Infof("Member Score: %+v", result)
	return &result, nil
}
