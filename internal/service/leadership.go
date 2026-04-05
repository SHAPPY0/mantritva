package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LeadershipService interface {
	Create(dto.LeadershipReq) (*model.Leadership, error)
}

type leadershipService struct {
	Repository            repository.LeadershipRepository
	PartyRepository       repository.PartyRepository
	StateRepository       repository.StateRepository
	DistrictRepository    repository.DistrictRepository
	LoksabhaRepository    repository.LoksabhaRepository
	VidhansabhaRepository repository.VidhansabhaRepository
}

func NewLeadershipService(repos *repository.Repos) LeadershipService {
	return &leadershipService{
		Repository:            repos.Leadership,
		PartyRepository:       repos.Party,
		StateRepository:       repos.State,
		DistrictRepository:    repos.District,
		LoksabhaRepository:    repos.Loksabha,
		VidhansabhaRepository: repos.VidhanSabha,
	}
}

func (s *leadershipService) Create(req dto.LeadershipReq) (*model.Leadership, error) {
	leadership := &model.Leadership{
		RulerID:       req.RulerID,
		MemberID:      req.MemberID,
		PartyID:       req.PartyID,
		Position:      utils.Position(req.Position),
		Role:          req.Role,
		Departments:   req.Departments,
		IsCentral:     req.IsCentral,
		StateID:       req.StateID,
		LoksabhaID:    req.LoksabhaID,
		VidhansabhaID: req.VidhansabhaID,
		RajyasabhaID:  req.RajyasabhaID,
		TermStart:     req.TermStart,
		IsActive:      req.IsActive,
	}
	if !req.TermEnd.IsNull() {
		leadership.TermEnd = req.TermEnd
	}
	if err := s.Repository.Create(leadership); err != nil {
		return nil, err
	}
	return leadership, nil
}
