package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type EntityService interface {
	List() ([]dto.EntityRes, error)
	Create(dto.EntityReq) (*model.Entity, error)
}

type entityService struct {
	Repository repository.EntityRepository
}

func NewEntityService(repos *repository.Repos) EntityService {
	return &entityService{
		Repository: repos.Entity,
	}
}

func (s *entityService) Create(req dto.EntityReq) (*model.Entity, error) {
	entity := &model.Entity{
		Name:  req.Name,
		Abbrv: req.Abbrv,
		// IsCentral:  *req.IsCentral,
		StateID:    req.StateID,
		DistrictID: req.DistrictID,
	}
	if err := s.Repository.Create(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *entityService) List() ([]dto.EntityRes, error) {
	entities, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return entities, nil
}
