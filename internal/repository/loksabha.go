package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type LoksabhaRepository interface {
	Create(*model.Loksabha) error
	List() ([]model.Loksabha, error)
	GetByID(id int) (*model.Loksabha, error)
	GetByStateID(stateId int) ([]model.Loksabha, error)
	GetByStateIDAndName(int, string) (*model.Loksabha, error)
}

type loksabhaRepository struct {
	DB *gorm.DB
}

func NewLoksabhaRepository(db *gorm.DB) LoksabhaRepository {
	return &loksabhaRepository{
		DB: db,
	}
}

func (r *loksabhaRepository) Create(loksabha *model.Loksabha) error {
	if err := r.DB.Create(loksabha).Error; err != nil {
		return err
	}
	return nil
}

func (r *loksabhaRepository) List() ([]model.Loksabha, error) {
	var loksabhas []model.Loksabha
	if err := r.DB.Find(&loksabhas).Error; err != nil {
		return nil, err
	}
	return loksabhas, nil
}

func (r *loksabhaRepository) GetByID(id int) (*model.Loksabha, error) {
	var loksabha model.Loksabha
	if err := r.DB.Where("id = ? AND is_deleted = false", id).First(&loksabha).Error; err != nil {
		return nil, err
	}
	return &loksabha, nil
}

func (r *loksabhaRepository) GetByStateID(stateId int) ([]model.Loksabha, error) {
	var loksabhas []model.Loksabha
	if err := r.DB.Where("state_id = ?", stateId).Order("name ASC").Find(&loksabhas).Error; err != nil {
		return nil, err
	}
	if len(loksabhas) == 0 {
		return nil, nil
	}
	return loksabhas, nil
}

func (r *loksabhaRepository) GetByStateIDAndName(stateId int, name string) (*model.Loksabha, error) {
	var loksabha model.Loksabha
	if err := r.DB.Where("state_id = ? AND name = ?", stateId, name).Find(&loksabha).Error; err != nil {
		return nil, err
	}
	if loksabha.ID <= 0 {
		return nil, nil
	}
	return &loksabha, nil
}
