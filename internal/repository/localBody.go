package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type LocalBodyRepository struct {
	DB *gorm.DB
}

func NewLocalBodyRepository(db *gorm.DB) LocalBodyRepository {
	return LocalBodyRepository{
		DB: db,
	}
}

func (r *LocalBodyRepository) IsExist(name string, districtId uint) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.LocalBody{}).Where("name = ? AND district_id = ? AND deleted_at IS NULL", name, districtId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LocalBodyRepository) GetByID(id uint) (*model.LocalBody, error) {
	var localBody model.LocalBody
	if err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&localBody).Error; err != nil {
		return nil, err
	}
	return &localBody, nil
}

func (r *LocalBodyRepository) CreateLocalBody(localBody *model.LocalBody) error {
	if err := r.DB.Create(localBody).Error; err != nil {
		return err
	}
	return nil
}

func (r *LocalBodyRepository) ListLocalBodiesByDistrict(districtId string) ([]model.LocalBody, error) {
	var localBodies []model.LocalBody
	if err := r.DB.Where("district_id = ?", districtId).Find(&localBodies).Error; err != nil {
		return nil, err
	}
	return localBodies, nil
}

func (r *LocalBodyRepository) ListLocalBodiesByStateAndDistrict(stateIds []uint, districtId string) ([]model.LocalBody, error) {
	var localBodies []model.LocalBody
	if err := r.DB.Where("district_id IN (select id from districts where state_id IN ? and id = ?)", stateIds, districtId).Find(&localBodies).Error; err != nil {
		return nil, err
	}
	return localBodies, nil
}

func (r *LocalBodyRepository) ListLocalBodiesByDistrictIDs(districtIds []uint) ([]model.LocalBody, error) {
	var localBodies []model.LocalBody
	if err := r.DB.Where("district_id IN ?", districtIds).Find(&localBodies).Error; err != nil {
		return nil, err
	}
	return localBodies, nil
}

func (r *LocalBodyRepository) ListLocalBodiesByIDAndDistrict(localBodyIDs []uint, districtId string) ([]model.LocalBody, error) {
	var localBodies []model.LocalBody
	if err := r.DB.Where("id IN ? AND district_id = ?", localBodyIDs, districtId).Find(&localBodies).Error; err != nil {
		return nil, err
	}
	return localBodies, nil
}
