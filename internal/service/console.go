package service

import (
	"log"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type ConsoleService interface {
	// Create(dto.NewsPressReq) (*model.NewsPressPost, error)
	GetUsersList() (*[]model.User, error)
	GetStatesList() (*[]model.State, error)
	GetPartiesList() (*[]model.Party, error)
	GetAlliancesList() (*[]model.Alliance, error)
	GetMembersList() (*[]dto.MemberRes, error)
	GetRulersList() (*[]dto.RulerRes, error)
	GetLeadershipByRulerID(int) (*[]dto.LeadershipRes, error)
	GetNewsPressList() (*[]dto.NewsPressPostRes, error)
}

type consoleService struct {
	UserRepository       repository.UserRepository
	StateRepository      repository.StateRepository
	PartyRepository      repository.PartyRepository
	AllianceRepository   repository.AllianceRepository
	MemberRepository     repository.MemberRepository
	RulerRepository      repository.RulerRepository
	LeadershipRepository repository.LeadershipRepository
	NewsPressRepository  repository.NewsPressPostRepository
}

func NewConsoleService(repos *repository.Repos) ConsoleService {
	return &consoleService{
		UserRepository:       repos.User,
		StateRepository:      repos.State,
		PartyRepository:      repos.Party,
		AllianceRepository:   repos.Alliance,
		MemberRepository:     repos.Member,
		RulerRepository:      repos.Ruler,
		LeadershipRepository: repos.Leadership,
		NewsPressRepository:  repos.NewsPressPost,
	}
}

func (s *consoleService) GetUsersList() (*[]model.User, error) {
	var users *[]model.User
	var err error
	users, err = s.UserRepository.ConsoleGetUsersList()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *consoleService) GetStatesList() (*[]model.State, error) {
	var states *[]model.State
	var err error
	states, err = s.StateRepository.ConsoleGetStateList()
	if err != nil {
		return nil, err
	}
	return states, nil
}

func (s *consoleService) GetPartiesList() (*[]model.Party, error) {
	var parties *[]model.Party
	var err error
	parties, err = s.PartyRepository.ConsoleGetPartiesList()
	if err != nil {
		return nil, err
	}
	return parties, nil
}

func (s *consoleService) GetAlliancesList() (*[]model.Alliance, error) {
	var alliances *[]model.Alliance
	var err error
	alliances, err = s.AllianceRepository.ConsoleGetAllianceList()
	if err != nil {
		return nil, err
	}
	return alliances, nil
}

func (s *consoleService) GetMembersList() (*[]dto.MemberRes, error) {
	var result []dto.MemberRes
	members, err := s.MemberRepository.ConsoleGetMembersList()
	if err != nil {
		return nil, err
	}
	for _, member := range *members {
		dtoMember := dto.MemberRes{Member: member}
		//Get Party
		party, err := s.PartyRepository.GetByMemberID(member.ID)
		if err != nil {
			log.Printf("Failed to get party details by memberID %d", member.ID)
		} else {
			dtoMember.Party = dto.MemberParty{
				Name:        party.Name,
				Abbrv:       party.Abbrv,
				Symbol:      party.Symbol,
				JoiningDate: party.JoiningDate,
				ExitDate:    party.ExitDate,
				Position:    party.Position,
			}
		}
		result = append(result, dtoMember)
	}
	return &result, nil
}

func (s *consoleService) GetRulersList() (*[]dto.RulerRes, error) {
	var rulersList []dto.RulerRes
	var rulers *[]model.Ruler
	var err error
	rulers, err = s.RulerRepository.ConsoleGetRulersList()
	if err != nil {
		return nil, err
	}
	for _, ruler := range *rulers {
		rulerItem := dto.RulerRes{Ruler: ruler}
		if ruler.StateID > 0 {
			state, err := s.StateRepository.GetByID(ruler.StateID)
			if err != nil {
				log.Printf("Error: While getting State detail %d, %v", ruler.StateID, err)
			} else {
				rulerItem.State = &model.StateDetail{State: *state}
			}
		}

		if ruler.PartyID > 0 {
			party, err := s.PartyRepository.GetByID(ruler.PartyID)
			if err != nil {
				log.Printf("Error: While getting Party detail %d, %v", ruler.PartyID, err)
			} else {
				rulerItem.Party = dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol}
			}
		}

		if ruler.IsAlliance {
			alliance, err := s.AllianceRepository.GetByID(ruler.AllianceID)
			if err != nil {
				log.Printf("Error: While getting Alliance detail %d, %v", ruler.AllianceID, err)
			} else {
				rulerItem.Alliance = dto.AllianceName{Name: alliance.Name, Abbrv: alliance.Abbrv}
			}
		}
		rulersList = append(rulersList, rulerItem)
	}
	return &rulersList, nil
}

func (s *consoleService) GetLeadershipByRulerID(rulerID int) (*[]dto.LeadershipRes, error) {
	var results []dto.LeadershipRes
	var leaderships []model.Leadership
	var err error
	leaderships, err = s.LeadershipRepository.GetByRulerID(uint(rulerID), true)
	if err != nil {
		return nil, err
	}
	for _, leadership := range leaderships {
		leader := dto.LeadershipRes{Leadership: leadership}
		if leader.MemberID > 0 {
			member, err := s.MemberRepository.GetByID(leader.MemberID)
			if err != nil {
				log.Printf("Error: While getting Member Detail %d, %v", leader.MemberID, err)
			} else {
				leader.Member = dto.Member{Name: member.Name, Photo: member.Photo}
			}
		}
		if leader.PartyID > 0 {
			party, err := s.PartyRepository.GetByID(int(leader.PartyID))
			if err != nil {
				log.Printf("Error: While getting Party Detail %d, %v", leader.PartyID, err)
			} else {
				leader.Party = dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol}
			}
		}
		if leader.StateID > 0 {
			state, err := s.StateRepository.GetByID(int(leader.StateID))
			if err != nil {
				log.Printf("Error: While getting State Detail %d, %v", leader.StateID, err)
			} else {
				leader.State = dto.State{Name: state.Name, Abbrv: state.Abbrv}
			}
		}
		results = append(results, leader)
	}

	return &results, nil
}

func (s *consoleService) GetNewsPressList() (*[]dto.NewsPressPostRes, error) {
	var results []dto.NewsPressPostRes
	var err error
	newsPressPosts, err := s.NewsPressRepository.ConsoleGetNewsPressList()
	if err != nil {
		return nil, err
	}
	for _, post := range newsPressPosts {
		newsPostRes := dto.NewsPressPostRes{NewsPressPost: post}
		if post.MemberId > 0 {
			member, err := s.MemberRepository.GetByID(post.MemberId)
			if err != nil {
				log.Printf("Error: While getting Member Detail %d, %v", post.MemberId, err)
			} else {
				newsPostRes.Member = dto.Member{Name: member.Name, Photo: member.Photo}
			}
		}
		if post.PartyId > 0 {
			party, err := s.PartyRepository.GetByID(int(post.PartyId))
			if err != nil {
				log.Printf("Error: While getting Party Detail %d, %v", post.PartyId, err)
			} else {
				newsPostRes.Party = dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol}
			}
		}
		results = append(results, newsPostRes)
	}
	return &results, nil
}
