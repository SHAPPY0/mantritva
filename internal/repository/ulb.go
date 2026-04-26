package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type UlbRepository struct {
	DB *gorm.DB
}

func NewULBRepository(db *gorm.DB) UlbRepository {
	return UlbRepository{
		DB: db,
	}
}

func (r *UlbRepository) CreateULB(ulb *model.ULB) error {
	if err := r.DB.Create(ulb).Error; err != nil {
		return err
	}
	return nil
}

func (r *UlbRepository) ListULBsByDistrict(districtId string) ([]model.ULB, error) {
	var ulbs []model.ULB
	if err := r.DB.Where("district_id = ? AND deleted_at IS NULL", districtId).Find(&ulbs).Error; err != nil {
		return nil, err
	}
	return ulbs, nil
}
