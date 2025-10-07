package service

import (
	"log"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type RulerService interface {
	List() ([]model.Ruler, error)
	Create(dto.RulerReq) (*model.Ruler, error)
	GetCentralRuler() (*dto.RulerRes, error)
	GetStateRuler(stateId string) (*dto.RulerRes, error)
	GetCenterLeadership() (map[string][]dto.Leader, error)
}

type rulerService struct {
	Repository           repository.RulerRepository
	PartyRepository      repository.PartyRepository
	AllianceRepository   repository.AllianceRepository
	LeadershipRepository repository.LeadershipRepository
	MemberRepository     repository.MemberRepository
	StateRepository      repository.StateRepository
}

func NewRulerService(repos *repository.Repos) RulerService {
	return &rulerService{
		Repository:           repos.Ruler,
		PartyRepository:      repos.Party,
		AllianceRepository:   repos.Alliance,
		LeadershipRepository: repos.Leadership,
		MemberRepository:     repos.Member,
		StateRepository:      repos.State,
	}
}

func (s *rulerService) Create(req dto.RulerReq) (*model.Ruler, error) {
	ruler := &model.Ruler{
		PartyID:    req.PartyID,
		IsAlliance: req.IsAlliance,
		AllianceID: req.AllianceID,
		IsCentral:  req.IsCentral,
		StateID:    req.StateID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		IsActive:   req.IsActive,
	}

	if err := s.Repository.Create(ruler); err != nil {
		return nil, err
	}
	return ruler, nil
}

func (s *rulerService) List() ([]model.Ruler, error) {
	rulers, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return rulers, nil
}

func (s *rulerService) GetCentralRuler() (*dto.RulerRes, error) {
	var rulerDetail dto.RulerRes
	ruler, err := s.Repository.GetCentralRuler()
	if err != nil {
		return nil, err
	}
	if ruler != nil {
		rulerDetail.Ruler = *ruler
		if ruler.PartyID > 0 {
			party, err := s.PartyRepository.GetByID(ruler.PartyID)
			if err != nil {
				log.Printf("Failed to get party details by partyID")
			} else {
				rulerDetail.Party = dto.Party{
					Name:   party.Name,
					Abbrv:  party.Abbrv,
					Symbol: party.Symbol,
				}
			}
		}
		if ruler.AllianceID > 0 {
			rulerDetail.IsAlliance = true
			alliance, err := s.AllianceRepository.GetByID(ruler.AllianceID)
			if err != nil {
				log.Printf("Failed to get alliance details by allianceID")
			} else {
				rulerDetail.Alliance = dto.AllianceName{
					Name:  alliance.Name,
					Abbrv: alliance.Abbrv,
				}
			}
		}
		if ruler.ID > 0 {
			leadership, err := s.LeadershipRepository.GetByRulerID(ruler.ID)
			if err != nil {
				log.Printf("Failed to get leadership details by rulerID")
			} else {
				for _, leader := range leadership {
					member, err := s.MemberRepository.GetByID(leader.MemberID)
					if err != nil {
						log.Printf("Failed to get member details by memberID %s", leader.MemberID)
					} else {
						rulerDetail.Leaders = append(rulerDetail.Leaders, dto.Leader{
							Name:        member.Name,
							Designation: dto.Position{Position: leader.Position, Role: leader.Role},
							Photo:       member.Photo,
							TermStart:   leader.TermStart,
							// TermEnd:   *leader.TermEnd,
						})
					}
				}
			}
		}
	}
	return &rulerDetail, nil
}

func (s *rulerService) GetStateRuler(stateId string) (*dto.RulerRes, error) {
	var rulerDetail dto.RulerRes
	ruler, err := s.Repository.GetStateRuler(stateId)
	if err != nil {
		return nil, err
	}
	if ruler != nil {
		rulerDetail.Ruler = *ruler
		//Fetch Party
		if ruler.PartyID > 0 {
			party, err := s.PartyRepository.GetByID(ruler.PartyID)
			if err != nil {
				log.Printf("Failed to get party details by partyID")
			} else {
				rulerDetail.Party = dto.Party{
					Name:   party.Name,
					Abbrv:  party.Abbrv,
					Symbol: party.Symbol,
				}
			}
		}
		//Fetch Alliance
		if ruler.AllianceID > 0 {
			rulerDetail.IsAlliance = true
			alliance, err := s.AllianceRepository.GetByID(ruler.AllianceID)
			if err != nil {
				log.Printf("Failed to get alliance details by allianceID")
			} else {
				rulerDetail.Alliance = dto.AllianceName{
					Name:  alliance.Name,
					Abbrv: alliance.Abbrv,
				}
			}
		}
		//Fetch Leaders
		if ruler.ID > 0 {
			leadership, err := s.LeadershipRepository.GetByRulerID(ruler.ID)
			if err != nil {
				log.Printf("Failed to get leadership details by rulerID")
			} else {
				for _, leader := range leadership {
					member, err := s.MemberRepository.GetByID(leader.MemberID)
					if err != nil {
						log.Printf("Failed to get member details by memberID %d", leader.MemberID)
					} else {
						rulerDetail.Leaders = append(rulerDetail.Leaders, dto.Leader{
							Name:        member.Name,
							Designation: dto.Position{Position: leader.Position, Role: leader.Role, Departments: leader.Departments},
							Photo:       member.Photo,
							TermStart:   leader.TermStart,
							// TermEnd:   *leader.TermEnd,
						})
					}
				}
			}
		}
	}
	return &rulerDetail, nil
}

func (s *rulerService) GetCenterLeadership() (map[string][]dto.Leader, error) {
	result := make(map[string][]dto.Leader)
	centerRuler, err := s.Repository.GetCentralRuler()
	if err != nil {
		return nil, err
	}
	if centerRuler != nil && centerRuler.ID > 0 {
		leadership, err := s.LeadershipRepository.GetByRulerIDAndPosition(centerRuler.ID, "Cabinet Minister")
		if err != nil {
			log.Printf("Failed to get leadership details by rulerID")
		} else {
			var leaders []dto.Leader
			for _, leader := range leadership {
				member, err := s.MemberRepository.GetByID(leader.MemberID)
				if err != nil {
					log.Printf("Failed to get member details by memberID %d", leader.MemberID)
				} else {
					cabMinLeader := dto.Leader{
						Name:        member.Name,
						Designation: dto.Position{Position: leader.Position, Role: leader.Role, Departments: leader.Departments},
						Photo:       member.Photo,
						TermStart:   leader.TermStart,
					}
					// if !leader.TermEnd.IsNull() {
					// 	cabMinLeader.TermEnd = *leader.TermEnd
					// }
					leaders = append(leaders, cabMinLeader)
				}
			}
			result["cabinet_minister"] = leaders
		}
	}
	return result, nil
}
