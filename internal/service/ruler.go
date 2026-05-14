package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type RulerService interface {
	List() ([]map[string]interface{}, error)
	Create(dto.RulerReq) (*model.Ruler, error)
	GetCentralRuler() (*dto.RulerRes, error)
	GetCenterLeadership() (map[string][]dto.Leader, error)
	GetStateRuler(string) (*dto.RulerRes, error)
	GetStateLeadership(string) ([]map[string]interface{}, error)
	GetRajyasabhaByStateID(string) ([]map[string]interface{}, error)
}

type rulerService struct {
	Repository           repository.RulerRepository
	PartyRepository      repository.PartyRepository
	AllianceRepository   repository.AllianceRepository
	LeadershipRepository repository.LeadershipRepository
	MemberRepository     repository.MemberRepository
	StateRepository      repository.StateRepository
	LoksabhRepository    repository.LoksabhaRepository
}

func NewRulerService(repos *repository.Repos) RulerService {
	return &rulerService{
		Repository:           repos.Ruler,
		PartyRepository:      repos.Party,
		AllianceRepository:   repos.Alliance,
		LeadershipRepository: repos.Leadership,
		MemberRepository:     repos.Member,
		StateRepository:      repos.State,
		LoksabhRepository:    repos.Loksabha,
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

func (s *rulerService) List() ([]map[string]interface{}, error) {
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
				log.Errorf("Failed to get party details by partyID %d", ruler.PartyID)
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
				log.Errorf("Failed to get alliance details by allianceID %d", ruler.AllianceID)
			} else {
				rulerDetail.Alliance = dto.AllianceName{
					Name:  alliance.Name,
					Abbrv: alliance.Abbrv,
				}
			}
		}
		if ruler.ID > 0 {
			positions := []string{string(utils.PRIME_MINISTER), string(utils.PRESIDENT), string(utils.VICE_PRESIDENT)}
			leadership, err := s.LeadershipRepository.GetByRulerIDAndPositions(ruler.ID, "0", positions, true)
			if err != nil {
				log.Errorf("Failed to get leadership details by rulerID %d", ruler.ID)
			} else {
				for _, leader := range leadership {
					var leadership dto.LeadershipRes
					err := utils.MapToStruct(leader, &leadership)
					if err != nil {
						log.Errorf("Failed to parse member data from leadership", nil)
						continue
					} else {
						utils.TransformPositionLabel(&leadership.Position.Position)
						rulerDetail.Leaders = append(rulerDetail.Leaders, dto.Leader{
							Name:        leadership.Member.Name,
							Designation: dto.Position{Position: leadership.Position.Position, Role: leadership.Position.Role, Departments: leadership.Position.Departments},
							Photo:       leadership.Member.Photo,
							TermStart:   leadership.TermStart,
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
				log.Errorf("Failed to get party details by partyID %d", ruler.PartyID)
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
				log.Errorf("Failed to get alliance details by allianceID %d", ruler.AllianceID)
			} else {
				rulerDetail.Alliance = dto.AllianceName{
					Name:  alliance.Name,
					Abbrv: alliance.Abbrv,
				}
			}
		}
		//Fetch Leaders
		if ruler.ID > 0 {
			leadership, err := s.LeadershipRepository.GetByRulerIDAndPosition(ruler.ID, string(utils.CHIEF_MINISTER), false)
			if err != nil {
				log.Errorf("Failed to get leadership details by rulerID %d", ruler.ID)
			} else {
				for _, leader := range leadership {
					member, err := s.MemberRepository.GetByID(leader.MemberID)
					if err != nil {
						log.Errorf("Failed to get member details by memberID %d", leader.MemberID)
					} else {
						rulerDetail.Leaders = append(rulerDetail.Leaders, dto.Leader{
							Name:        member.Name,
							Designation: dto.Position{Position: leader.Position.Label(), Role: leader.Role, Departments: leader.Departments},
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

func (s *rulerService) GetStateLeadership(stateId string) ([]map[string]interface{}, error) {
	// result := make(map[string][]dto.Leader)
	stateRuler, err := s.Repository.GetStateRuler(stateId)
	if err != nil {
		return nil, err
	}
	if stateRuler != nil && stateRuler.ID > 0 {
		postitons := []string{string(utils.CHIEF_MINISTER), string(utils.MLA), string(utils.DEPUTY_CHIEF_MINISTER)}
		leaders, err := s.LeadershipRepository.GetByRulerIDAndPositions(stateRuler.ID, stateId, postitons, false)
		if err != nil {
			log.Errorf("Failed to get state leadership by rulerID %d", stateRuler.ID)
		}
		return leaders, nil
	}
	return nil, err
}

func (s *rulerService) GetCenterLeadership() (map[string][]dto.Leader, error) {
	result := make(map[string][]dto.Leader)
	centerRuler, err := s.Repository.GetCentralRuler()
	if err != nil {
		return nil, err
	}
	if centerRuler != nil && centerRuler.ID > 0 {
		leadership, err := s.LeadershipRepository.GetByRulerIDAndPosition(centerRuler.ID, string(utils.CABINET_MINISTER), true)
		if err != nil {
			log.Errorf("Failed to get leadership details by rulerID %d", centerRuler.ID)
		} else {
			var leaders []dto.Leader
			for _, leader := range leadership {
				var member dto.Member
				memberInfo, err := s.MemberRepository.GetMemberById(leader.MemberID)
				if err != nil {
					log.Errorf("Error fetching member info by id %d", leader.MemberID)
				} else {
					if err := utils.MapToStruct(memberInfo, &member); err != nil {
						log.Errorf("Error converting member map to dto.member")
						continue
					} else {
						utils.TransformPositionLabel(&member.Position)
						item := dto.Leader{
							ID:          leader.MemberID,
							Name:        member.Name,
							Designation: dto.Position{Position: leader.Position.Label(), Role: leader.Role, Departments: leader.Departments},
							Photo:       member.Photo,
							TermStart:   leader.TermStart,
						}
						if !utils.IsEmpty(member.Party) {
							item.Party = dto.Party{
								Name:   member.Party.Name,
								Abbrv:  member.Party.Abbrv,
								Symbol: member.Party.Symbol,
							}
						}
						if leader.StateID > 0 {
							state, err := s.StateRepository.GetByID(int(leader.StateID))
							if err == nil {
								item.State = state.Name
							}
						}
						if leader.LoksabhaID > 0 {
							loksabha, err := s.LoksabhRepository.GetByID(int(leader.LoksabhaID))
							if err == nil {
								item.Loksabha = loksabha.Name
							}
						}
						leaders = append(leaders, item)
					}
				}
			}
			result["cabinet_minister"] = leaders
		}

		mpLeadership, err := s.LeadershipRepository.GetByRulerIDAndPosition(centerRuler.ID, string(utils.MP_LS), true)
		if err != nil {
			log.Errorf("Failed to get MPs details by rulerID %d", centerRuler.ID)
		} else {
			var leaders []dto.Leader
			for _, leader := range mpLeadership {
				var member dto.Member
				memberInfo, err := s.MemberRepository.GetMemberById(leader.MemberID)
				if err != nil {
					log.Errorf("Failed to get member details by memberID %d", leader.MemberID)
				} else {
					if err := utils.MapToStruct(memberInfo, &member); err != nil {
						log.Errorf("Error converting member map to dto.member")
						continue
					} else {
						item := dto.Leader{
							ID:          leader.MemberID,
							Name:        member.Name,
							Designation: dto.Position{Position: leader.Position.Label(), Role: leader.Role, Departments: leader.Departments},
							Photo:       member.Photo,
							TermStart:   leader.TermStart,
						}
						if !utils.IsEmpty(member.Party) {
							item.Party = dto.Party{
								Name:   member.Party.Name,
								Abbrv:  member.Party.Abbrv,
								Symbol: member.Party.Symbol,
							}
						}
						if leader.StateID > 0 {
							state, err := s.StateRepository.GetByID(int(leader.StateID))
							if err == nil {
								item.State = state.Name
							}
						}
						if leader.LoksabhaID > 0 {
							loksabha, err := s.LoksabhRepository.GetByID(int(leader.LoksabhaID))
							if err == nil {
								item.Loksabha = loksabha.Name
							}
						}
						leaders = append(leaders, item)
					}
				}
			}
			result["member_of_parliament"] = leaders
		}

		mosLeadership, err := s.LeadershipRepository.GetByRulerIDAndPosition(centerRuler.ID, string(utils.MINISTER_OF_STATE), true)
		if err != nil {
			log.Errorf("Failed to get MOS details by rulerID %d", centerRuler.ID)
		} else {
			var leaders []dto.Leader
			for _, leader := range mosLeadership {
				var member dto.Member
				// member, err := s.MemberRepository.GetByID(leader.MemberID)
				memberInfo, err := s.MemberRepository.GetMemberById(leader.MemberID)
				if err != nil {
					log.Errorf("Failed to get member details by memberID %d", leader.MemberID)
				} else {
					if err := utils.MapToStruct(memberInfo, &member); err != nil {
						log.Errorf("Error converting member map to dto.member")
						continue
					} else {
						item := dto.Leader{
							Name:        member.Name,
							Designation: dto.Position{Position: leader.Position.Label(), Role: leader.Role, Departments: leader.Departments},
							Photo:       member.Photo,
							TermStart:   leader.TermStart,
						}
						if !utils.IsEmpty(member.Party) {
							item.Party = dto.Party{
								Name:   member.Party.Name,
								Abbrv:  member.Party.Abbrv,
								Symbol: member.Party.Symbol,
							}
						}
						if leader.StateID > 0 {
							state, err := s.StateRepository.GetByID(int(leader.StateID))
							if err == nil {
								item.State = state.Name
							}
						}
						if leader.LoksabhaID > 0 {
							loksabha, err := s.LoksabhRepository.GetByID(int(leader.LoksabhaID))
							log.Errorf("Loksabha fetch error: %v", err)
							if err == nil {
								item.Loksabha = loksabha.Name
							}
						}

						leaders = append(leaders, item)
					}
				}
			}
			result["member_of_state"] = leaders
		}
	}
	return result, nil
}

func (s *rulerService) GetRajyasabhaByStateID(stateId string) ([]map[string]interface{}, error) {
	rsMembers, err := s.LeadershipRepository.GetRajyaSabhMembersByStateID(stateId)
	if err != nil {
		return nil, err
	}
	return rsMembers, nil
}
