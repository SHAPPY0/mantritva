package service

import (

	// "path/filepath"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
)

type MemberService interface {
	List() ([]dto.MemberRes, error)
	Create(c *gin.Context, req dto.MemberReq, photo *multipart.FileHeader) (*model.Member, error)
	BulkMemberUpload(c *gin.Context, filePath string) (bool, error)
	GetByName(name string) (*dto.MemberRes, error)
	ListUnaffiliatedMembers() ([]dto.MemberRes, error)
	GetCareerByMemberID(int) ([]map[string]interface{}, error)
	GetPostsByMemberID(dto.Pagination, int) (*[]dto.NewsPressPostRes, error)
	GetGrievancesByMemberID(dto.Pagination, int, string) (*[]dto.GrievanceShort, error)
	// GetByID(id uint) (*dto.Member, error)
}

type memberService struct {
	Repository           repository.MemberRepository
	LeadershipRepository repository.LeadershipRepository
	PartyRepository      repository.PartyRepository
	NewsPressRepository  repository.NewsPressPostRepository
	GrievanceRepository  repository.GrievanceRepository
	ReportRepository     repository.ReportRepository
}

func NewMemberService(repos *repository.Repos) MemberService {
	return &memberService{
		Repository:           repos.Member,
		LeadershipRepository: repos.Leadership,
		PartyRepository:      repos.Party,
		NewsPressRepository:  repos.NewsPressPost,
		GrievanceRepository:  repos.Grievance,
		ReportRepository:     repos.Report,
	}
}

func (s *memberService) Create(ctx *gin.Context, req dto.MemberReq, photo *multipart.FileHeader) (*model.Member, error) {
	destinationPath := filepath.Join("members")
	storage := storage.NewLocalStorage(config.Env.GetString("FILE_BASE_PATH"), destinationPath)
	fileMeta, err := storage.Save(ctx, photo)
	if err != nil {
		return nil, err
	}
	member := &model.Member{
		Name:             req.Name,
		FullName:         req.FullName,
		DOB:              req.DOB,
		Gender:           req.Gender,
		Qualifications:   req.Qualifications,
		MaritalStatus:    req.MaritalStatus,
		FinancialAssets:  req.FinancialAssets,
		SocialMediaLinks: req.SocialMediaLinks,
		Photo:            fileMeta.AbsolutePath,
	}
	if err := s.Repository.Create(member); err != nil {
		storage.Delete(fileMeta.FullPath)
		return nil, err
	}
	if len(req.CriminalRecords) > 0 {
		for _, record := range req.CriminalRecords {
			criminalRecord := &model.CriminalRecord{
				MemberID:       member.ID,
				Offense:        record.Offense,
				Status:         record.Status,
				ConvictionDate: record.ConvictionDate,
				Sentence:       record.Sentence,
			}
			if err := s.Repository.CreateCriminalRecord(criminalRecord); err != nil {
				return nil, err
			}
		}
	}

	return member, nil
}

func (s *memberService) List() ([]dto.MemberRes, error) {
	var result []dto.MemberRes
	members, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	for _, member := range members {
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
	return result, nil
}

func (s *memberService) BulkMemberUpload(c *gin.Context, filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode JSON into struct
	var memberReq dto.MemberReq
	if err := json.NewDecoder(file).Decode(&memberReq); err != nil {
		panic(err)
	}
	return true, nil
}

func (s *memberService) GetByName(name string) (*dto.MemberRes, error) {
	var data dto.MemberRes
	memberName := strings.Replace(name, "_", " ", -1)
	member, err := s.Repository.GetByName(memberName)
	if err != nil {
		return nil, err
	}
	data.Member = *member
	//Get Party
	party, err := s.PartyRepository.GetByMemberID(member.ID)
	if err != nil {
		return nil, err
	}
	data.Party = *party

	//Get position from leadership table
	leadership, err := s.LeadershipRepository.GetByMemberID(member.ID)
	if err != nil {
		return nil, err
	}
	data.Position = leadership.Position
	data.Departments = leadership.Departments

	score, err := s.ReportRepository.GetMemberScore(member.ID, time.Now().Format("2006"))
	if err != nil {
		log.Warnf("No score found for memberID %d", member.ID)
	} else if score != nil {
		data.Score = &dto.MemberScore{
			Score:     score.TotalScore,
			MaxScore:  score.MaxScore,
			Year:      score.Year,
			UpdatedAt: score.CreatedAt.String(),
		}
	}

	return &data, nil
}

func (s *memberService) GetCareerByMemberID(memberID int) ([]map[string]interface{}, error) {
	careers, err := s.LeadershipRepository.GetCareerByMemberID(memberID)
	if err != nil {
		return nil, err
	}
	return careers, nil
}

func (s *memberService) ListUnaffiliatedMembers() ([]dto.MemberRes, error) {
	var result []dto.MemberRes
	members, err := s.Repository.ListUnaffiliatedMembers()
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		dtoMember := dto.MemberRes{Member: member}
		result = append(result, dtoMember)
	}
	return result, nil
}

func (s *memberService) GetPostsByMemberID(pagination dto.Pagination, memberID int) (*[]dto.NewsPressPostRes, error) {
	var newsPressPostRes []dto.NewsPressPostRes
	posts, err := s.NewsPressRepository.GetPostsByMemberID(pagination, memberID)
	if err != nil {
		return nil, err
	}
	for _, post := range *posts {
		dtoPost := dto.NewsPressPostRes{NewsPressPost: post}
		memberAndParty, err := s.Repository.GetMemberAndPartyDetails(memberID, int(post.PartyId))
		if err != nil {
			return nil, err
		}
		if err := utils.MapToStruct(memberAndParty, &dtoPost); err != nil {
			return nil, err
		}
		newsPressPostRes = append(newsPressPostRes, dtoPost)
	}
	return &newsPressPostRes, nil
}

func (s *memberService) GetGrievancesByMemberID(pagination dto.Pagination, memberID int, position string) (*[]dto.GrievanceShort, error) {
	var grievancesRes []dto.GrievanceShort
	grievances, err := s.GrievanceRepository.GetGrievancesByMemberID(pagination, memberID, position)
	if err != nil {
		return nil, err
	}
	for _, grievance := range *grievances {
		dtoGrievance := dto.GrievanceShort{
			GrievanceID:    grievance.GrievanceID,
			Title:          grievance.Title,
			Type:           grievance.Type,
			SubType:        grievance.SubType,
			Level:          grievance.Level,
			CenterSubLevel: grievance.CenterSubLevel,
			DateOfIncident: utils.DateOnly{Time: &grievance.DateOfIncident},
			Location:       grievance.Location,
			Impact:         grievance.Impact,
			FinalStatus:    grievance.FinalStatus,
			CreatedAt:      grievance.CreatedAt,
		}
		grievancesRes = append(grievancesRes, dtoGrievance)
	}
	return &grievancesRes, nil
}
