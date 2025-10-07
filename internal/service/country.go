package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type CountryService interface {
	List() ([]model.Country, error)
	Create(dto.CountryReq) (*model.Country, error)
}

type countryService struct {
	Repository repository.CountryRepository
}

func NewCountryService(repo repository.CountryRepository) CountryService {
	return &countryService{
		Repository: repo,
	}
}

func (s *countryService) Create(req dto.CountryReq) (*model.Country, error) {
	country := &model.Country{
		Name:    req.Name,
		Abbrv:   req.Abbrv,
		Icon:    req.Icon,
		Capital: req.Capital,
	}

	if err := s.Repository.Create(country); err != nil {
		return nil, err
	}
	return country, nil
}

func (s *countryService) List() ([]model.Country, error) {
	countries, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return countries, nil
}
