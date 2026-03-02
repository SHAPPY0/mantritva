package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type EnquiryRepository interface {
	Create(*model.Enquiry) error
	GetAll() ([]*model.Enquiry, error)
}

type enquiryRepository struct {
	DB *gorm.DB
}

func NewEnquiryRepository(db *gorm.DB) EnquiryRepository {
	return &enquiryRepository{
		DB: db,
	}
}

func (r *enquiryRepository) Create(data *model.Enquiry) error {
	if err := r.DB.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *enquiryRepository) GetAll() ([]*model.Enquiry, error) {
	var enquiries []*model.Enquiry
	if err := r.DB.Find(&enquiries).Order("created_at DESC").Error; err != nil {
		return nil, err
	}
	return enquiries, nil
}
