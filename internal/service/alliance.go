package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type AllianceService interface {
	List() ([]model.Alliance, error)
	Create(dto.AllianceReq) (*model.Alliance, error)
	AddPartyToAlliance(dto.AlliancePartyReq) error
	GetAllianceByAbbrv(abbrv string) (*dto.AlliancePartyDetail, error)
}

type allianceService struct {
	Repository repository.AllianceRepository
}

func NewAllianceService(repo repository.AllianceRepository) AllianceService {
	return &allianceService{
		Repository: repo,
	}
}

func (s *allianceService) Create(req dto.AllianceReq) (*model.Alliance, error) {
	alliance := &model.Alliance{
		Name:       req.Name,
		Abbrv:      req.Abbrv,
		IsActive:   req.IsActive,
		FormedYear: req.FormedYear,
		Remark:     req.Remark,
	}
	if err := s.Repository.Create(alliance); err != nil {
		return nil, err
	}
	return alliance, nil
}

func (s *allianceService) List() ([]model.Alliance, error) {
	alliances, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return alliances, nil
}

func (s *allianceService) AddPartyToAlliance(req dto.AlliancePartyReq) error {
	allianceParty := &model.AllianceParty{
		AllianceID:  req.AllianceID,
		PartyID:     req.PartyID,
		IsLeader:    req.IsLeader,
		JoiningDate: req.JoiningDate,
		LeavingDate: req.LeavingDate,
		IsActive:    req.IsActive,
	}
	if err := s.Repository.AddPartyToAlliance(allianceParty); err != nil {
		return err
	}
	return nil
}

func (s *allianceService) GetAllianceByAbbrv(abbrv string) (*dto.AlliancePartyDetail, error) {
	alliance, err := s.Repository.GetByAbbrv(abbrv)
	if err != nil {
		return nil, err
	}
	return alliance, nil
}
