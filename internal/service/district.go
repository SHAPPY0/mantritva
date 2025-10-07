package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type DistrictService interface {
	List() ([]model.District, error)
	Create(dto.DistrictReq) (*model.District, error)
	GetByState(stateId string) ([]model.District, error)
}

type districtService struct {
	Repository repository.DistrictRepository
}

func NewDistrictService(repo repository.DistrictRepository) DistrictService {
	return &districtService{
		Repository: repo,
	}
}

func (s *districtService) Create(req dto.DistrictReq) (*model.District, error) {
	district := &model.District{
		Name:    req.Name,
		Abbrv:   req.Abbrv,
		StateID: req.StateID,
	}
	if err := s.Repository.Create(district); err != nil {
		return nil, err
	}
	return district, nil
}

func (s *districtService) List() ([]model.District, error) {
	states, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return states, nil
}

func (s *districtService) GetByState(stateId string) ([]model.District, error) {
	districts, err := s.Repository.GetByState(stateId)
	if err != nil {
		return nil, err
	}
	if len(districts) == 0 {
		return nil, nil // No districts found for the given state
	}
	return districts, nil
}
