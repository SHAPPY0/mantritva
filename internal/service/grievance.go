package service

import (
	"crypto/rand"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
)

type GrievanceSerivce interface {
	Create(*gin.Context, dto.Grievance) (*model.Grievance, error)
	ListByUser(string) ([]model.Grievance, error)
	GetByGrievanceID(string) (*dto.GrievanceRes, error)
	GetByStateID(int) ([]*dto.GrievanceShort, error)
}

type grievanceSerivce struct {
	Repository   repository.GrievanceRepository
	MemberRepo   repository.MemberRepository
	PartyRepo    repository.PartyRepository
	StateRepo    repository.StateRepository
	DistrictRepo repository.DistrictRepository
	UserRepo     repository.UserRepository
	EntityRepo   repository.EntityRepository
}

func NewGrievanceService(repos *repository.Repos) GrievanceSerivce {
	return &grievanceSerivce{
		Repository:   repos.Grievance,
		MemberRepo:   repos.Member,
		PartyRepo:    repos.Party,
		StateRepo:    repos.State,
		DistrictRepo: repos.District,
		UserRepo:     repos.User,
		EntityRepo:   repos.Entity,
	}
}

func (s *grievanceSerivce) Create(ctx *gin.Context, req dto.Grievance) (*model.Grievance, error) {
	g := &model.Grievance{
		GrievanceID:          generateID(5),
		Type:                 req.Type,
		SubType:              req.SubType,
		Jurisdiction:         req.Jurisdiction,
		CenterSubLevel:       req.CenterSubLevel,
		CabinetMinisterID:    req.CabinetMinisterID,
		MemberOfParliamentID: req.MemberOfParliamentID,
		EntityID:             req.EntityID,
		StateID:              req.StateID,
		DistrictID:           req.DistrictID,
		StateLeaderID:        req.StateLeaderID,
		Title:                req.Title,
		Description:          req.Description,
		HideIdentity:         req.HideIdentity,
		DateOfIncident:       *req.DateOfIncident.Time,
		Location:             req.Location,
		Severity:             req.Severity,
		Impact:               req.Impact,
		SLA:                  req.SLA.Time,
		FinalStatus:          utils.UNDERREVIEW,
		Evidences:            make([]model.Document, 0),
		StatusRevisions:      make([]model.StatusRevisions, 0),
		AuditFieldsBSON:      model.InsertBSONAuditFields(req.CreatedBy),
	}

	destinationPath := filepath.Join("grievances", g.GrievanceID)
	storage, err := storage.New(destinationPath)
	if err != nil {
		log.Infof("Error initiating storage connection %v", err)
		return nil, err
	}

	for _, evidenceReq := range req.Evidences {
		fileMeta, err := storage.Save(ctx, evidenceReq.File)
		if err != nil {
			log.Warnf("Error while storing file %v", err)
		} else {
			g.Evidences = append(g.Evidences, model.Document{
				Url:      fileMeta.DownloadURL,
				FileName: fileMeta.FullPath,
				FileType: fileMeta.Extenstion,
				Remark:   evidenceReq.Comment,
			})
		}
	}
	if err := s.Repository.Create(g); err != nil {
		storage.DeleteAll(g.Evidences)
		return nil, err
	}
	return g, nil
}

func generateID(n int) string {
	if n <= 0 {
		return ""
	}
	var sb strings.Builder
	// first digit 1-9 (avoid leading zero)
	first, _ := rand.Int(rand.Reader, big.NewInt(9))
	sb.WriteByte(byte('1' + first.Int64()))
	for i := 1; i < (n - 1); i++ {
		d, _ := rand.Int(rand.Reader, big.NewInt(10))
		sb.WriteByte(byte('0' + d.Int64()))
	}
	year := time.Now().Format("06")
	return "mg" + year + sb.String()
}

func (s *grievanceSerivce) ListByUser(userID string) ([]model.Grievance, error) {
	grievances, err := s.Repository.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	return grievances, nil
}

func (s *grievanceSerivce) GetByGrievanceID(grievanceID string) (*dto.GrievanceRes, error) {
	var grievanceRes dto.GrievanceRes
	grievance, err := s.Repository.GetByGrievanceID(grievanceID)
	if err != nil {
		return nil, err
	}

	grievanceRes.Grievance = *grievance
	if grievance.CabinetMinisterID != 0 {
		member, err := s.MemberRepo.GetMemberById(grievance.CabinetMinisterID)
		if err != nil {
			log.Error("Error fetching member details by ID")
		}
		var memberDto dto.Member
		if err := utils.MapToStruct(member, &memberDto); err != nil {
			log.Error("Error converting member map to dto.member")
		}
		grievanceRes.CabinetMinister = &memberDto
	}
	if grievance.MemberOfParliamentID != 0 {
		member, err := s.MemberRepo.GetMemberById(grievance.MemberOfParliamentID)
		if err != nil {
			log.Error("Error fetching member details by ID")
		}
		var memberDto dto.Member
		if err := utils.MapToStruct(member, &memberDto); err != nil {
			log.Error("Error converting member map to dto.member")
		}
		grievanceRes.MemberOfParliament = &memberDto
	}
	if grievance.StateLeaderID != 0 {
		member, err := s.MemberRepo.GetMemberById(grievance.StateLeaderID)
		if err != nil {
			log.Error("Error fetching member details by ID")
		}
		var memberDto dto.Member
		if err := utils.MapToStruct(member, &memberDto); err != nil {
			log.Error("Error converting member map to dto.member")
		}
		grievanceRes.StateLeader = &memberDto
	}
	if grievance.StateID != 0 {
		state, err := s.StateRepo.GetByID(int(grievance.StateID))
		if err != nil {
			log.Error("Error fetching state details by ID")
		} else {
			grievanceRes.State = &model.State{Name: state.Name}
		}
	}
	if grievance.DistrictID != 0 {
		district, err := s.DistrictRepo.GetByID(grievance.DistrictID)
		if err != nil {
			log.Error("Error fetching district details by ID")
		} else {
			grievanceRes.District = district.Name
		}
	}
	if grievance.EntityID != 0 {
		ok, entity, err := s.EntityRepo.GetByID(int(grievance.EntityID))
		if err != nil {
			log.Error("Error fetching district details by ID")
		} else if ok && entity != nil {
			grievanceRes.Entity = &dto.Entity{
				Name:         entity.Name,
				Abbrv:        entity.Abbrv,
				IsCentral:    entity.IsCentral,
				StateName:    entity.StateName,
				DistrictName: entity.DistrictName,
			}
		}
	}

	if !grievance.HideIdentity && grievance.CreatedBy != "" {
		//convert string to int
		id, err := utils.StrToInt(grievance.CreatedBy)
		if err != nil {
			log.Error("Error converting CreatedBy to int")
		}
		user, err := s.UserRepo.GetByID(id)
		if err != nil {
			log.Errorf("Error fetching user details by ID %d", id)
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

func (s *grievanceSerivce) GetByStateID(stateID int) ([]*dto.GrievanceShort, error) {
	grievancesList := make([]*dto.GrievanceShort, 0)
	grievances, err := s.Repository.GetByStateID(stateID)
	if err != nil {
		return nil, err
	}
	for _, g := range *grievances {
		grievancesList = append(grievancesList, &dto.GrievanceShort{
			GrievanceID:    g.GrievanceID,
			Type:           g.Type,
			SubType:        g.SubType,
			Jurisdiction:   g.Jurisdiction,
			Title:          g.Title,
			DateOfIncident: utils.DateOnly{Time: &g.DateOfIncident},
			Location:       g.Location,
			Severity:       string(g.Severity),
			Impact:         g.Impact,
			FinalStatus:    g.FinalStatus,
			CreatedAt:      g.CreatedAt,
			UpvoteCount:    g.UpvoteCount,
			DownvoteCount:  g.DownvoteCount,
		})
	}
	return grievancesList, nil
}
