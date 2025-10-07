package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type StateService interface {
	List() ([]model.State, error)
	Create(dto.StateReq) (*model.State, error)
	Get(name string) (*model.StateDetail, error)
}

type stateService struct {
	Repository repository.StateRepository
}

func NewStateService(repo repository.StateRepository) StateService {
	return &stateService{
		Repository: repo,
	}
}

func (s *stateService) Create(req dto.StateReq) (*model.State, error) {
	state := &model.State{
		Name:      req.Name,
		Abbrv:     req.Abbrv,
		CountryID: req.CountryID,
		Capital:   req.Capital,
		Structure: req.Structure,
	}

	if err := s.Repository.Create(state); err != nil {
		return nil, err
	}
	return state, nil
}

func (s *stateService) List() ([]model.State, error) {
	states, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return states, nil
}

func (s *stateService) Get(name string) (*model.StateDetail, error) {
	state, err := s.Repository.Get(name)
	if err != nil {
		return nil, err
	}
	return state, nil
}
