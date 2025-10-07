package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type CountryRepository interface {
	Create(*model.Country) error
	List() ([]model.Country, error)
}

type countryRepository struct {
	DB *gorm.DB
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{
		DB: db,
	}
}

func (r *countryRepository) Create(country *model.Country) error {
	if err := r.DB.Create(country).Error; err != nil {
		return err
	}
	return nil
}

func (r *countryRepository) List() ([]model.Country, error) {
	var countries []model.Country
	if err := r.DB.Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}
