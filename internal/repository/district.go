package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type DistrictRepository interface {
	Create(*model.District) error
	List() ([]model.District, error)
	GetByState(string) ([]model.District, error)
	GetByID(uint) (*model.District, error)
	ConsoleGetDistrictList() (*[]model.District, error)
	ConsoleGetDistrictListByStateIDs([]uint) (*[]model.District, error)
	ConsoleGetDistrictListByID([]uint) (*[]model.District, error)
}

type districtRepository struct {
	DB *gorm.DB
}

func NewDistrictRepository(db *gorm.DB) DistrictRepository {
	return &districtRepository{
		DB: db,
	}
}

func (r *districtRepository) Create(district *model.District) error {
	if err := r.DB.Create(district).Error; err != nil {
		return err
	}
	return nil
}

func (r *districtRepository) List() ([]model.District, error) {
	var districts []model.District
	if err := r.DB.Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

func (r *districtRepository) GetByState(stateId string) ([]model.District, error) {
	var districts []model.District
	if err := r.DB.Where("state_id = ?", stateId).Order("name ASC").Find(&districts).Error; err != nil {
		return nil, err
	}
	if len(districts) == 0 {
		return nil, nil
	}
	return districts, nil
}

func (r *districtRepository) GetByID(id uint) (*model.District, error) {
	var district model.District
	if err := r.DB.Where("id = ?", id).Find(&district).Error; err != nil {
		return nil, err
	}
	return &district, nil
}

func (r *districtRepository) ConsoleGetDistrictList() (*[]model.District, error) {
	var districts []model.District
	if err := r.DB.Find(&districts).Error; err != nil {
		return nil, err
	}
	return &districts, nil
}

func (r *districtRepository) ConsoleGetDistrictListByStateIDs(stateIDs []uint) (*[]model.District, error) {
	var districts []model.District
	if err := r.DB.Where("state_id IN ?", stateIDs).Find(&districts).Error; err != nil {
		return nil, err
	}
	return &districts, nil
}

func (r *districtRepository) ConsoleGetDistrictListByID(districtIDs []uint) (*[]model.District, error) {
	var districts []model.District
	if err := r.DB.Where("id IN ?", districtIDs).Find(&districts).Error; err != nil {
		return nil, err
	}
	return &districts, nil
}
