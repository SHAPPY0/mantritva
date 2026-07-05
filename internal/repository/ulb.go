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

func (r *UlbRepository) IsExist(name string, districtId uint) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.ULB{}).Where("name = ? AND district_id = ? AND deleted_at IS NULL", name, districtId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UlbRepository) GetByID(id uint) (*model.ULB, error) {
	var ulb model.ULB
	if err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&ulb).Error; err != nil {
		return nil, err
	}
	return &ulb, nil
}

func (r *UlbRepository) CreateULB(ulb *model.ULB) error {
	if err := r.DB.Create(ulb).Error; err != nil {
		return err
	}
	return nil
}

func (r *UlbRepository) ListULBsByDistrict(districtId string) ([]model.ULB, error) {
	var ulbs []model.ULB
	if err := r.DB.Where("district_id = ?", districtId).Find(&ulbs).Error; err != nil {
		return nil, err
	}
	return ulbs, nil
}

func (r *UlbRepository) ListULBsByStateAndDistrict(stateIds []uint, districtId string) ([]model.ULB, error) {
	var ulbs []model.ULB
	if err := r.DB.Where("district_id IN (select id from districts where state_id IN ? and id = ?)", stateIds, districtId).Find(&ulbs).Error; err != nil {
		return nil, err
	}
	return ulbs, nil
}

func (r *UlbRepository) ListLocalBodiesByDistrictIDs(districtIds []uint) ([]model.ULB, error) {
	var ulbs []model.ULB
	if err := r.DB.Where("district_id IN ?", districtIds).Find(&ulbs).Error; err != nil {
		return nil, err
	}
	return ulbs, nil
}

func (r *UlbRepository) ListLocalBodiesByIDAndDistrict(localBodyIDs []uint, districtId string) ([]model.ULB, error) {
	var ulbs []model.ULB
	if err := r.DB.Where("id IN ? AND district_id = ?", localBodyIDs, districtId).Find(&ulbs).Error; err != nil {
		return nil, err
	}
	return ulbs, nil
}
