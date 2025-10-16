package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type LoksabhaService interface {
	List() ([]model.Loksabha, error)
	Create(dto.LoksabhaReq) (*model.Loksabha, error)
	GetByStateID(stateId int) ([]model.Loksabha, error)
}

type loksabhaService struct {
	Repository repository.LoksabhaRepository
}

func NewLoksabhaService(repo repository.LoksabhaRepository) LoksabhaService {
	return &loksabhaService{
		Repository: repo,
	}
}

func (s *loksabhaService) Create(req dto.LoksabhaReq) (*model.Loksabha, error) {
	loksabha := &model.Loksabha{
		Name:       req.Name,
		StateID:    req.StateID,
		DistrictID: req.DistrictID,
		PCNumber:   req.PCNumber,
	}
	if err := s.Repository.Create(loksabha); err != nil {
		return nil, err
	}
	return loksabha, nil
}

func (s *loksabhaService) List() ([]model.Loksabha, error) {
	loksabhas, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return loksabhas, nil
}

func (s *loksabhaService) GetByStateID(stateId int) ([]model.Loksabha, error) {
	loksabhas, err := s.Repository.GetByStateID(stateId)
	if err != nil {
		return nil, err
	}
	return loksabhas, nil
}
