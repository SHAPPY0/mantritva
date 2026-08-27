package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsoleService interface {
	// Create(dto.NewsPressReq) (*model.NewsPressPost, error)
	GetUsersList() ([]dto.UserRes, error)
	GetStatesList() (*[]model.State, error)
	GetStatesByUser(auth.AuthUser) (*[]model.State, error)
	GetDistrictsByStateByUser(auth.AuthUser) (*[]model.District, error)
	GetLocalBodiesByDistrictByUser(auth.AuthUser, string) (*[]model.ULB, error)
	GetStatesByID([]uint) (*[]model.State, error)
	GetPartiesList() (*[]model.Party, error)
	GetAlliancesList() (*[]model.Alliance, error)
	GetMembersList() (*[]dto.MemberRes, error)
	GetByID(id uint) (*dto.MemberRes, error)
	GetRulersList() (*[]dto.RulerRes, error)
	GetRulersListByUser(auth.AuthUser) (*[]dto.RulerRes, error)
	GetCurrentRulersList() ([]map[string]interface{}, error)
	GetLeadershipByRulerID(int, bool) ([]map[string]interface{}, error)
	GetNewsPressList() (*[]dto.NewsPressPostRes, error)
	GetNewsPressPostByID(string) (*model.NewsPressPost, error)
	UpdateNewsPressPost(string, map[string]interface{}) (bool, error)
	GetGrievancesList() (*[]model.Grievance, error)
	GetGrievanceByID(string) (*dto.GrievanceRes, error)
	GrievanceUpdateStatus(string, dto.GrievanceStatus) (bool, error)
	GetDevelopmentWelfaresList() (*[]model.DevelopmentWelfare, error)
	GetDevelopmentWelfareByID(string) (*dto.DevelopmentWelfareRes, error)
	DevelopmentWelfareUpdateStatus(string, dto.DevelopmentWelfareStatus) (bool, error)
	ExitLeadership(dto.LeadershipExitReq) (bool, error)
	UpdateMemberByID(*gin.Context, uint, map[string]interface{}, *multipart.FileHeader) (bool, error)
	UpdateRole(dto.UpdateUserReq) (bool, error)
	GetEnquiriesList() ([]*model.Enquiry, error)
	AddRole(string, dto.AddRoleReq) (*model.UserRole, error)
	GetUserRoles(string) ([]model.UserRole, error)
	SetupBulkLeader(string, multipart.File) (map[string]string, map[string]string, error)
	GetDashboardStats() (*dto.ConsoleDashboardRes, error)
}

type consoleService struct {
	UserRepository        repository.UserRepository
	UserRoleRepository    repository.UserRoleRepository
	StateRepository       repository.StateRepository
	PartyRepository       repository.PartyRepository
	AllianceRepository    repository.AllianceRepository
	MemberRepository      repository.MemberRepository
	RulerRepository       repository.RulerRepository
	LeadershipRepository  repository.LeadershipRepository
	NewsPressRepository   repository.NewsPressPostRepository
	EntityRepository      repository.EntityRepository
	GrievanceRepository           repository.GrievanceRepository
	DevelopmentWelfareRepository  repository.DevelopmentWelfareRepository
	DistrictRepository            repository.DistrictRepository
	ReportRepository      repository.ReportRepository
	EnquiryRepository     repository.EnquiryRepository
	LoksabhaRepository    repository.LoksabhaRepository
	VidhansabhaRepository repository.VidhansabhaRepository
	LocalBodyRepository   repository.UlbRepository
}

func NewConsoleService(repos *repository.Repos) ConsoleService {
	return &consoleService{
		UserRepository:        repos.User,
		UserRoleRepository:    repos.UserRole,
		StateRepository:       repos.State,
		PartyRepository:       repos.Party,
		AllianceRepository:    repos.Alliance,
		MemberRepository:      repos.Member,
		RulerRepository:       repos.Ruler,
		LeadershipRepository:  repos.Leadership,
		NewsPressRepository:   repos.NewsPressPost,
		EntityRepository:      repos.Entity,
		GrievanceRepository:          repos.Grievance,
		DevelopmentWelfareRepository: repos.DevelopmentWelfare,
		DistrictRepository:           repos.District,
		ReportRepository:      repos.Report,
		EnquiryRepository:     repos.Enquiry,
		LoksabhaRepository:    repos.Loksabha,
		VidhansabhaRepository: repos.VidhanSabha,
		LocalBodyRepository:   repos.Ulb,
	}
}

func (s *consoleService) GetUsersList() ([]dto.UserRes, error) {
	result := make([]dto.UserRes, 0)
	var users *[]model.User
	var err error
	users, err = s.UserRepository.ConsoleGetUsersList()
	if err != nil {
		return nil, err
	}
	var userIds []uint
	for _, user := range *users {
		userIds = append(userIds, user.ID)
	}
	roles, err := s.UserRoleRepository.GetRoleByUserIDs(userIds)
	if err != nil {
		return nil, err
	}
	roleMap := make(map[uint][]model.UserRole)
	for _, role := range roles {
		roleMap[role.UserID] = append(roleMap[role.UserID], role)
	}
	for _, user := range *users {
		var userRes dto.UserRes
		userRes.User = user
		userRes.Roles = make([]dto.RoleRes, 0)
		if userRoles, ok := roleMap[user.ID]; ok {
			for _, role := range userRoles {
				roleRes := dto.RoleRes{
					Role:    role.Role,
					Scope:   role.Scope,
					ScopeID: role.ScopeID,
				}
				if roleRes.Scope == utils.ROLE_SCOPE_STATE {
					state, err := s.StateRepository.GetByID(int(role.ScopeID))
					if err != nil {
						log.Errorf("Error while getting State detail %d, %v", role.ScopeID, err)
					} else if state != nil {
						roleRes.State = state.Name
					}
				} else if roleRes.Scope == utils.ROLE_SCOPE_DISTRICT {
					district, err := s.DistrictRepository.GetByID(role.ScopeID)
					if err != nil {
						log.Errorf("Error while getting District detail %d, %v", role.ScopeID, err)
					} else if district != nil {
						roleRes.District = district.Name
					}
				} else if roleRes.Scope == utils.ROLE_SCOPE_LOKSABHA {
					loksabha, err := s.LoksabhaRepository.GetByID(int(role.ScopeID))
					if err != nil {
						log.Errorf("Error while getting Loksabha detail %d, %v", role.ScopeID, err)
					} else if loksabha != nil {
						roleRes.Loksabha = loksabha.Name
					}
				} else if roleRes.Scope == utils.ROLE_SCOPE_VIDHANSABHA {
					vidhansabha, err := s.VidhansabhaRepository.GetByID(int(role.ScopeID))
					if err != nil {
						log.Errorf("Error while getting Vidhansabha detail %d, %v", role.ScopeID, err)
					} else if vidhansabha != nil {
						roleRes.Vidhansabha = vidhansabha.Name
					}
				} else if roleRes.Scope == utils.ROLE_SCOPE_LOCALBODY {
					localBody, err := s.LocalBodyRepository.GetByID(role.ScopeID)
					if err != nil {
						log.Errorf("Error while getting LocalBody detail %d, %v", role.ScopeID, err)
					} else if localBody != nil {
						roleRes.LocalBody = localBody.Name
					}
				} else if roleRes.Scope == utils.ROLE_SCOPE_WARD {
					// roleRes.Ward, _ = s.ReportRepository.GetWardByID(int(role.ScopeID))
				}
				userRes.Roles = append(userRes.Roles, roleRes)
			}
		}
		result = append(result, userRes)
	}
	return result, nil
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

func (s *consoleService) GetStatesByUser(user auth.AuthUser) (*[]model.State, error) {
	var states *[]model.State
	var err error
	if user.IsAdminOrSuperAdmin() {
		states, err = s.StateRepository.ConsoleGetStateList()
	} else if user.IsStateAdmin() {
		stateIDs := user.GetAssignedStates()
		states, err = s.StateRepository.ConsoleGetStateListByID(stateIDs)
	} else if user.IsDistrictAdmin() {
		districtIDs := user.GetAssignedDistricts()
		states, err = s.StateRepository.ConsoleGetStateListByDistrictIDs(districtIDs)
	}
	if err != nil {
		return nil, err
	}
	return states, nil
}

func (s *consoleService) GetDistrictsByStateByUser(user auth.AuthUser) (*[]model.District, error) {
	var districts *[]model.District
	var err error
	if user.IsAdminOrSuperAdmin() {
		districts, err = s.DistrictRepository.ConsoleGetDistrictList()
	} else if user.IsStateAdmin() {
		stateIDs := user.GetAssignedStates()
		districts, err = s.DistrictRepository.ConsoleGetDistrictListByStateIDs(stateIDs)
	} else if user.IsDistrictAdmin() {
		districtIDs := user.GetAssignedDistricts()
		districts, err = s.DistrictRepository.ConsoleGetDistrictListByID(districtIDs)
	}
	if err != nil {
		return nil, err
	}
	return districts, nil
}

func (s *consoleService) GetLocalBodiesByDistrictByUser(user auth.AuthUser, districtID string) (*[]model.ULB, error) {
	var localBodies []model.ULB
	var err error
	if user.IsAdminOrSuperAdmin() {
		localBodies, err = s.LocalBodyRepository.ListULBsByDistrict(districtID)
	} else if user.IsStateAdmin() {
		stateIDs := user.GetAssignedStates()
		localBodies, err = s.LocalBodyRepository.ListULBsByStateAndDistrict(stateIDs, districtID)
	} else if user.IsDistrictAdmin() {
		districtIDs := user.GetAssignedDistricts()
		if !utils.Contains(districtIDs, utils.StrToUint(districtID)) {
			return nil, utils.NewError("Unauthorized access to district %s", districtID)
		}
		localBodies, err = s.LocalBodyRepository.ListULBsByDistrict(districtID)
	} else if user.IsLocalBodyAdmin() {
		localBodyIDs := user.GetAssignedLocalBodies()
		localBodies, err = s.LocalBodyRepository.ListLocalBodiesByIDAndDistrict(localBodyIDs, districtID)
	}
	if err != nil {
		return nil, err
	}
	return &localBodies, nil
}

func (s *consoleService) GetStatesByID(stateIDs []uint) (*[]model.State, error) {
	var states *[]model.State
	var err error
	states, err = s.StateRepository.ConsoleGetStateListByID(stateIDs)
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
	var rulers *[]model.Ruler
	var err error
	rulers, err = s.RulerRepository.ConsoleGetRulersList()
	if err != nil {
		return nil, err
	}
	rulersRes, err := s.PrepareRulerRes(rulers)
	if err != nil {
		return nil, err
	}
	return rulersRes, nil
}

func (s *consoleService) GetRulersListByUser(user auth.AuthUser) (*[]dto.RulerRes, error) {
	var rulers *[]model.Ruler
	var err error
	if user.IsAdminOrSuperAdmin() {
		rulers, err = s.RulerRepository.ConsoleGetRulersList()
	} else if user.IsStateAdmin() {
		stateIDs := user.GetAssignedStates()
		rulers, err = s.RulerRepository.ConsoleGetRulersListByStateIDs(stateIDs)
	} else if user.IsDistrictAdmin() {
		districtIDs := user.GetAssignedDistricts()
		rulers, err = s.RulerRepository.ConsoleGetRulersListByDistrictIDs(districtIDs)
	}
	if err != nil {
		return nil, err
	}
	return s.PrepareRulerRes(rulers)
}

func (s *consoleService) PrepareRulerRes(rulers *[]model.Ruler) (*[]dto.RulerRes, error) {
	var rulersList []dto.RulerRes
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
		if post.MemberID > 0 {
			member, err := s.MemberRepository.GetByID(post.MemberID)
			if err != nil {
				log.Errorf("Error: While getting Member Detail %d, %v", post.MemberID, err)
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
			ok, entity, err := s.EntityRepository.GetByID(int(post.EntityId))
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
		if post.CreatedBy != "" {
			createdByInt, err := utils.StrToInt(post.CreatedBy)
			if err != nil {
				log.Errorf("Error converting created by to int %s", post.CreatedBy)
			} else if createdByInt > 0 {
				user, err := s.UserRepository.GetByID(createdByInt)
				if err != nil {
					log.Errorf("Error getting user by id %d, %v", createdByInt, err)
				} else {
					newsPostRes.PostAdmin = &dto.PostAdmin{Name: user.Name, UserID: fmt.Sprintf("%d", user.ID)}
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
		} else {
			utils.TransformPositionLabel(&memberDto.Position)
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
		} else {
			utils.TransformPositionLabel(&memberDto.Position)
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
		} else {
			utils.TransformPositionLabel(&memberDto.Position)
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

func (s *consoleService) GetDevelopmentWelfaresList() (*[]model.DevelopmentWelfare, error) {
	developments, err := s.DevelopmentWelfareRepository.ConsoleGetDevelopmentWelfaresList()
	if err != nil {
		return nil, err
	}
	return &developments, nil
}

func (s *consoleService) GetDevelopmentWelfareByID(id string) (*dto.DevelopmentWelfareRes, error) {
	development, err := s.DevelopmentWelfareRepository.ConsoleGetDevelopmentWelfareByID(id)
	if err != nil {
		return nil, err
	}

	result := dto.DevelopmentWelfareRes{
		DevelopmentWelfare: *development,
	}

	if development.MemberID != 0 {
		member, err := s.MemberRepository.GetMemberById(development.MemberID)
		if err != nil {
			log.Errorf("Error fetching member details by ID %d", development.MemberID)
		} else {
			var memberDto dto.Member
			if err := utils.MapToStruct(member, &memberDto); err != nil {
				log.Errorf("Error converting member map to dto.member")
			} else {
				utils.TransformPositionLabel(&memberDto.Position)
				result.Member = memberDto
			}
		}
	}

	if development.PartyID != 0 {
		party, err := s.PartyRepository.GetByID(int(development.PartyID))
		if err != nil {
			log.Errorf("Error fetching party details by ID %d", development.PartyID)
		} else {
			result.Party = dto.Party{
				Name:   party.Name,
				Abbrv:  party.Abbrv,
				Symbol: party.Symbol,
			}
		}
	}

	if development.StateID != nil && *development.StateID != 0 {
		state, err := s.StateRepository.GetByID(int(*development.StateID))
		if err != nil {
			log.Errorf("Error fetching state details by ID %d", *development.StateID)
		} else {
			result.State = &model.State{Name: state.Name}
		}
	}

	if development.DistrictID != nil && *development.DistrictID != 0 {
		district, err := s.DistrictRepository.GetByID(*development.DistrictID)
		if err != nil {
			log.Errorf("Error fetching district details by ID %d", *development.DistrictID)
		} else {
			result.District = district.Name
		}
	}

	return &result, nil
}

func (s *consoleService) DevelopmentWelfareUpdateStatus(id string, status dto.DevelopmentWelfareStatus) (bool, error) {
	err := s.DevelopmentWelfareRepository.ConsoleUpdateDevelopmentWelfareStatus(id, status)
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

func (s *consoleService) GetEnquiriesList() ([]*model.Enquiry, error) {
	var enquiries []*model.Enquiry
	var err error
	enquiries, err = s.EnquiryRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return enquiries, nil
}

func (s *consoleService) AddRole(createdBy string, req dto.AddRoleReq) (*model.UserRole, error) {
	isExist, err := s.UserRoleRepository.IsRoleExist(req.UserID, req.RoleID, req.Scope, &req.ScopeID)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, utils.NewError("Role already assigned for user %d", req.UserID)
	}
	role, err := s.UserRoleRepository.AddRole(createdBy, req)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *consoleService) GetUserRoles(userID string) ([]model.UserRole, error) {
	roles, err := s.UserRoleRepository.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *consoleService) SetupBulkLeader(username string, file multipart.File) (map[string]string, map[string]string, error) {
	var leaderReqs []dto.BulkLeader
	if err := json.NewDecoder(file).Decode(&leaderReqs); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for i, leader := range leaderReqs {
		if leader.Member.Name == "" {
			log.Errorf("No valid member info at position %d", i)
			continue
		}
		//Member
		memberName := leader.Member.Name
		if memberName != "" {
			memberName = strings.ToUpper(string(leader.Member.Name[0])) + leader.Member.Name[1:]
		}
		memberModel := &model.Member{
			Name:             memberName,
			FullName:         leader.Member.FullName,
			DOB:              leader.Member.DOB,
			Gender:           leader.Member.Gender,
			Qualifications:   leader.Member.Qualifications,
			MaritalStatus:    leader.Member.MaritalStatus,
			FinancialAssets:  leader.Member.FinancialAssets,
			SocialMediaLinks: leader.Member.SocialMediaLinks,
			Photo:            "",
			AuditFields: model.AuditFields{
				CreatedBy: username,
			},
		}
		if err := s.MemberRepository.Create(memberModel); err != nil {
			createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error at position %d Error: %s", i, err.Error())).Error()
			continue
		}
		member, err := s.MemberRepository.GetByName(memberName)
		if err != nil {
			createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error at position %d Error: %s", i, err.Error())).Error()
			continue
		}

		var party *model.Party

		// Party
		if leader.Party.Setup {
			party, err = s.PartyRepository.GetByAbbrv(leader.Party.Name)
			if err != nil {
				createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error member %s  while getting party info: %s", memberName, err.Error())).Error()
				continue
			}
			if party == nil {
				createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("No party info found by name %s", leader.Party.Name)).Error()
				continue
			}

			partyMember := &model.PartyMembers{
				PartyID:     party.ID,
				MemberID:    member.ID,
				Position:    leader.Party.Position,
				JoiningDate: leader.Party.JoiningDate,
				ExitDate:    leader.Party.ExitDate,
				AuditFields: model.AuditFields{
					CreatedBy: username,
				},
			}
			if err := s.PartyRepository.AddMember(partyMember); err != nil {
				createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("No party info found by name %s", leader.Party.Name)).Error()
				continue
			}
		}

		//=======Leadership========

		if leader.Party.Setup && leader.Leadership.Setup {
			// Get Ruler
			var ruler *model.Ruler

			if leader.Leadership.IsCentral {
				ruler, err = s.RulerRepository.GetCentralRuler()
				if err != nil {
					log.Errorf("Error while getting ruler info %v", err)
					createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error setting up leadership for member %s, %v", leader.Member.Name, err)).Error()
					continue
				}
			} else if leader.Leadership.State != "" {
				stateRow, err := s.StateRepository.Get(leader.Leadership.State)
				if err != nil {
					log.Errorf("Error while getting state info %v", err)
					createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error setting up leadership for member %s, %v", leader.Member.Name, err)).Error()
					continue
				} else if stateRow != nil {
					ruler, err = s.RulerRepository.GetStateRuler(fmt.Sprintf("%d", stateRow.ID))
					if err != nil {
						log.Errorf("Error while getting state info %v", err)
						createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error setting up leadership for member %s, %v", leader.Member.Name, err)).Error()
						continue
					}
				} else {
					log.Errorf("No state info found %s", leader.Leadership.State)
					createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error setting up leadership for member %s, %v", leader.Member.Name, err)).Error()
					continue
				}

			}

			leadership := &model.Leadership{
				RulerID:     ruler.ID,
				MemberID:    member.ID,
				PartyID:     party.ID,
				Position:    utils.Position(leader.Leadership.Position),
				Role:        leader.Leadership.Role,
				Departments: leader.Leadership.Departments,
				IsCentral:   leader.Leadership.IsCentral,
				StateID:     uint(ruler.StateID),
				// LoksabhaID:    req.LoksabhaID,
				// VidhansabhaID: req.VidhansabhaID,
				// RajyasabhaID:  req.RajyasabhaID,
				TermStart: leader.Leadership.TermStart,
				IsActive:  true,
				AuditFields: model.AuditFields{
					CreatedBy: username,
				},
			}
			log.Infof("leadership:: %v", leadership)
			// if !leader.Leadership.TermEnd.IsNull() || !leader.Leadership.TermEnd.IsZero() {
			// 	leadership.TermEnd = leader.Leadership.TermEnd
			// }
			if err := s.LeadershipRepository.Create(leadership); err != nil {
				createError[leader.Member.Name] = utils.NewError(fmt.Sprintf("Error setting up leadership for member %s, %v", leader.Member.Name, err)).Error()
			} else {
				createSuccess[leader.Member.Name] = "Leader setup successfully"
			}
		} else {
			createSuccess[leader.Member.Name] = "Leader setup successfully"
		}

	}
	return createSuccess, createError, nil
}

func (s *consoleService) GetDashboardStats() (*dto.ConsoleDashboardRes, error) {
	users, err := s.UserRepository.ConsoleCountUsers()
	if err != nil {
		return nil, err
	}
	newsPress, err := s.NewsPressRepository.ConsoleCountNewsPress()
	if err != nil {
		return nil, err
	}
	developmentWelfares, err := s.DevelopmentWelfareRepository.ConsoleCountDevelopmentWelfares()
	if err != nil {
		return nil, err
	}
	grievances, err := s.GrievanceRepository.ConsoleCountGrievances()
	if err != nil {
		return nil, err
	}
	return &dto.ConsoleDashboardRes{
		Users:               users,
		NewsPress:           newsPress,
		DevelopmentWelfares: developmentWelfares,
		Grievances:          grievances,
		Total:               users + newsPress + developmentWelfares + grievances,
	}, nil
}
