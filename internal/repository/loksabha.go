package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type LoksabhaRepository interface {
	Create(*model.Loksabha) error
	List() ([]model.Loksabha, error)
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
