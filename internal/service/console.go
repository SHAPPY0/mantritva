package service

import (
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsoleService interface {
	// Create(dto.NewsPressReq) (*model.NewsPressPost, error)
	GetUsersList() (*[]model.User, error)
	GetStatesList() (*[]model.State, error)
	GetPartiesList() (*[]model.Party, error)
	GetAlliancesList() (*[]model.Alliance, error)
	GetMembersList() (*[]dto.MemberRes, error)
	GetByID(id uint) (*dto.MemberRes, error)
	GetRulersList() (*[]dto.RulerRes, error)
	GetCurrentRulersList() ([]map[string]interface{}, error)
	GetLeadershipByRulerID(int, bool) ([]map[string]interface{}, error)
	GetNewsPressList() (*[]dto.NewsPressPostRes, error)
	GetNewsPressPostByID(string) (*model.NewsPressPost, error)
	UpdateNewsPressPost(string, map[string]interface{}) (bool, error)
	GetGrievancesList() (*[]model.Grievance, error)
	GetGrievanceByID(string) (*dto.GrievanceRes, error)
	GrievanceUpdateStatus(string, dto.GrievanceStatus) (bool, error)
	ExitLeadership(dto.LeadershipExitReq) (bool, error)
	UpdateMemberByID(*gin.Context, uint, map[string]interface{}, *multipart.FileHeader) (bool, error)
	UpdateRole(dto.UpdateUserReq) (bool, error)
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
	EntityRepository     repository.EntityRepository
	GrievanceRepository  repository.GrievanceRepository
	DistrictRepository   repository.DistrictRepository
	ReportRepository     repository.ReportRepository
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
		EntityRepository:     repos.Entity,
		GrievanceRepository:  repos.Grievance,
		DistrictRepository:   repos.District,
		ReportRepository:     repos.Report,
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
			log.Errorf("Failed to get party details by memberID %d", member.ID)
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

func (s *consoleService) GetByID(id uint) (*dto.MemberRes, error) {
	var data dto.MemberRes
	member, err := s.MemberRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	data.Member = *member
	//Get Party
	// party, err := s.PartyRepository.GetByMemberID(member.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// data.Party = *party

	//Get position from leadership table
	// leadership, err := s.LeadershipRepository.GetByMemberID(member.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// data.Position = leadership.Position
	// data.Departments = leadership.Departments

	// score, err := s.ReportRepository.GetMemberScore(member.ID, time.Now().Format("2006"))
	// if err != nil {
	// 	log.Warnf("No score found for memberID %d", member.ID)
	// } else if score != nil {
	// 	data.Score = &dto.MemberScore{
	// 		Score:     score.TotalScore,
	// 		MaxScore:  score.MaxScore,
	// 		Year:      score.Year,
	// 		UpdatedAt: score.CreatedAt.String(),
	// 	}
	// }

	return &data, nil
}

func (s *consoleService) UpdateMemberByID(ctx *gin.Context, id uint, data map[string]interface{}, photo *multipart.FileHeader) (bool, error) {
	if photo != nil {
		destinationPath := filepath.Join("members")
		storage, err := storage.New(destinationPath)
		if err != nil {
			return false, err
		}
		fileMeta, err := storage.Save(ctx, photo)
		if err != nil {
			return false, err
		}
		data["photo"] = fileMeta.DownloadURL
	}
	ok, err := s.MemberRepository.UpdateMemberByID(id, data)
	if err != nil {
		return false, err
	}
	return ok, nil
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
				log.Errorf("Error while getting State detail %d, %v", ruler.StateID, err)
			} else {
				rulerItem.State = &model.StateDetail{State: *state}
			}
		}

		if ruler.PartyID > 0 {
			party, err := s.PartyRepository.GetByID(ruler.PartyID)
			if err != nil {
				log.Errorf("Error while getting Party detail %d, %v", ruler.PartyID, err)
			} else {
				rulerItem.Party = dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol}
			}
		}

		if ruler.IsAlliance {
			alliance, err := s.AllianceRepository.GetByID(ruler.AllianceID)
			if err != nil {
				log.Errorf("Error while getting Alliance detail %d, %v", ruler.AllianceID, err)
			} else {
				rulerItem.Alliance = dto.AllianceName{Name: alliance.Name, Abbrv: alliance.Abbrv}
			}
		}
		rulersList = append(rulersList, rulerItem)
	}
	return &rulersList, nil
}

func (s *consoleService) GetLeadershipByRulerID(rulerID int, isCentral bool) ([]map[string]interface{}, error) {
	// var results []dto.LeadershipRes
	// var leaderships []model.Leadership
	leaderships, err := s.LeadershipRepository.GetByRulerID(uint(rulerID), isCentral)
	if err != nil {
		return nil, err
	}

	return leaderships, nil
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
				log.Errorf("Error: While getting Member Detail %d, %v", post.MemberId, err)
			} else {
				newsPostRes.Member = dto.Member{Name: member.Name, Photo: member.Photo}
			}
		}
		if post.PartyId > 0 {
			party, err := s.PartyRepository.GetByID(int(post.PartyId))
			if err != nil {
				log.Errorf("Error: While getting Party Detail %d, %v", post.PartyId, err)
			} else {
				newsPostRes.Party = dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol}
			}
		}
		if post.EntityId > 0 {
			ok, entity, err := s.EntityRepository.GetEntityByID(int(post.EntityId))
			if err != nil {
				log.Errorf("Error: While getting Entity Detail %d, %v", post.EntityId, err)
			} else if ok {
				newsPostRes.Entity = dto.Entity{
					Name:         entity.Name,
					Abbrv:        entity.Abbrv,
					IsCentral:    entity.IsCentral,
					StateName:    entity.StateName,
					DistrictName: entity.DistrictName,
				}
			}
		}
		results = append(results, newsPostRes)
	}
	return &results, nil
}

func (s *consoleService) GetNewsPressPostByID(id string) (*model.NewsPressPost, error) {
	var err error
	var post *model.NewsPressPost
	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err = s.NewsPressRepository.ConsoleGetNewsPressPostByID(postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *consoleService) UpdateNewsPressPost(id string, data map[string]interface{}) (bool, error) {
	ok, err := s.NewsPressRepository.UpdatePost(id, data)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *consoleService) GetGrievancesList() (*[]model.Grievance, error) {
	var err error
	grievances, err := s.GrievanceRepository.ConsoleGetGrievancesList()
	if err != nil {
		return nil, err
	}
	return &grievances, nil
}

func (s *consoleService) GetGrievanceByID(grievanceID string) (*dto.GrievanceRes, error) {
	var err error
	var grievanceRes dto.GrievanceRes
	grievance, err := s.GrievanceRepository.ConsoleGetGrievanceByID(grievanceID)
	if err != nil {
		return nil, err
	}
	grievanceRes.Grievance = *grievance
	if grievance.CabinetMinisterID != 0 {
		member, err := s.MemberRepository.GetMemberById(grievance.CabinetMinisterID)
		if err != nil {
			log.Errorf("Error: Error fetching member details by ID")
		}
		var memberDto dto.Member
		if err := utils.MapToStruct(member, &memberDto); err != nil {
			log.Errorf("Error: Error converting member map to dto.member")
		}
		grievanceRes.CabinetMinister = &memberDto
	}
	if grievance.MemberOfParliamentID != 0 {
		member, err := s.MemberRepository.GetMemberById(grievance.MemberOfParliamentID)
		if err != nil {
			log.Errorf("Error: Error fetching member details by ID %d", grievance.MemberOfParliamentID)
		}
		var memberDto dto.Member
		if err := utils.MapToStruct(member, &memberDto); err != nil {
			log.Errorf("Error: Error converting member map to dto.member")
		}
		grievanceRes.MemberOfParliament = &memberDto
	}
	if grievance.StateLeaderID != 0 {
		member, err := s.MemberRepository.GetMemberById(grievance.StateLeaderID)
		if err != nil {
			log.Errorf("Error: Error fetching member details by ID")
		}
		var memberDto dto.Member
		if err := utils.MapToStruct(member, &memberDto); err != nil {
			log.Errorf("Error: Error converting member map to dto.member %+v", err)
		}
		grievanceRes.StateLeader = &memberDto
	}
	if grievance.StateID != 0 {
		state, err := s.StateRepository.GetByID(int(grievance.StateID))
		if err != nil {
			log.Errorf("Error: Error fetching state details by ID %d", grievance.StateID)
		} else {
			grievanceRes.State = &model.State{Name: state.Name}
		}
	}
	if grievance.DistrictID != 0 {
		district, err := s.DistrictRepository.GetByID(grievance.DistrictID)
		if err != nil {
			log.Errorf("Error: Error fetching district details by ID %d", grievance.DistrictID)
		} else {
			grievanceRes.District = district.Name
		}
	}
	if !grievance.HideIdentity && grievance.CreatedBy != "" {
		id, err := utils.StrToInt(grievance.CreatedBy)
		if err != nil {
			log.Errorf("Error: Error converting CreatedBy to int")
		}
		user, err := s.UserRepository.GetByID(id)
		if err != nil {
			log.Errorf("Error: Error fetching user details by ID %+v", err)
		} else {
			grievanceRes.User = dto.UserClaim{
				UserID: user.UserID,
				Name:   user.Name,
				Email:  user.Email,
			}
		}
	}
	return &grievanceRes, nil
}

func (s *consoleService) GrievanceUpdateStatus(grievanceID string, status dto.GrievanceStatus) (bool, error) {
	var err error
	err = s.GrievanceRepository.ConsoleUpdateGrievanceStatus(grievanceID, status)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *consoleService) GetCurrentRulersList() ([]map[string]interface{}, error) {
	rulers, err := s.RulerRepository.List()
	if err != nil {
		return nil, err
	}
	return rulers, nil
}

func (s *consoleService) ExitLeadership(req dto.LeadershipExitReq) (bool, error) {
	err := s.LeadershipRepository.ExitLeadership(req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *consoleService) UpdateRole(user dto.UpdateUserReq) (bool, error) {
	ok, err := s.UserRepository.UpdateRole(user)
	if err != nil {
		return false, err
	}
	return ok, nil
}
