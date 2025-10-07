package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type LeadershipService interface {
	Create(dto.LeadershipReq) (*model.Leadership, error)
}

type leadershipService struct {
	Repository      repository.LeadershipRepository
	PartyRepository repository.PartyRepository
}

func NewLeadershipService(repos *repository.Repos) LeadershipService {
	return &leadershipService{
		Repository:      repos.Leadership,
		PartyRepository: repos.Party,
	}
}

func (s *leadershipService) Create(req dto.LeadershipReq) (*model.Leadership, error) {
	leadership := &model.Leadership{
		RulerID:     req.RulerID,
		MemberID:    req.MemberID,
		PartyID:     req.PartyID,
		Position:    req.Position,
		Role:        req.Role,
		Departments: req.Departments,
		IsCentral:   req.IsCentral,
		StateID:     req.StateID,
		TermStart:   req.TermStart,
		TermEnd:     &req.TermEnd,
		IsActive:    req.IsActive,
	}
	if err := s.Repository.Create(leadership); err != nil {
		return nil, err
	}
	return leadership, nil
}
