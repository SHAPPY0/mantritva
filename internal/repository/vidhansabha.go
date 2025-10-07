package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type VidhansabhaRepository interface {
	Create(*model.Vidhansabha) error
	List() ([]model.Vidhansabha, error)
}

type vidhansabhaRepository struct {
	DB *gorm.DB
}

func NewVidhansabhaRepository(db *gorm.DB) VidhansabhaRepository {
	return &vidhansabhaRepository{
		DB: db,
	}
}

func (r *vidhansabhaRepository) Create(vidhanSabha *model.Vidhansabha) error {
	if err := r.DB.Create(vidhanSabha).Error; err != nil {
		return err
	}
	return nil
}

func (r *vidhansabhaRepository) List() ([]model.Vidhansabha, error) {
	var vidhanSabhas []model.Vidhansabha
	if err := r.DB.Find(&vidhanSabhas).Error; err != nil {
		return nil, err
	}
	return vidhanSabhas, nil
}
