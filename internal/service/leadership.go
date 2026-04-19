package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LeadershipService interface {
	Create(dto.LeadershipReq) (*model.Leadership, error)
	CreateBulk(file multipart.File) (map[string]string, map[string]string, error)
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
	if !leadership.IsCentral {
		exist, err := s.Repository.LeaderExist(leadership.RulerID, leadership.MemberID, leadership.StateID, string(leadership.Position))
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, utils.NewError("Leadership already exist")
		}
	}

	if err := s.Repository.Create(leadership); err != nil {
		return nil, err
	}
	return leadership, nil
}

func (s *leadershipService) CreateBulk(file multipart.File) (map[string]string, map[string]string, error) {
	var leadersReq []dto.LeadershipReq
	if err := json.NewDecoder(file).Decode(&leadersReq); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for _, leader := range leadersReq {
		memberID := fmt.Sprintf("%d", leader.MemberID)
		mem, err := s.Create(leader)
		if err != nil {
			createError[memberID] = err.Error()
		} else if mem != nil {
			createSuccess[string(memberID)] = "Success"
		}
	}
	return createSuccess, createError, nil
}
