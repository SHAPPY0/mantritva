package service

import (

	// "path/filepath"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/storage"
)

type MemberService interface {
	List() ([]dto.MemberRes, error)
	Create(c *gin.Context, req dto.MemberReq, photo *multipart.FileHeader) (*model.Member, error)
	BulkMemberUpload(c *gin.Context, filePath string) (bool, error)
	GetByName(name string) (*dto.MemberRes, error)
	GetCareerByMemberID(memberID int) ([]map[string]interface{}, error)
}

type memberService struct {
	Repository           repository.MemberRepository
	LeadershipRepository repository.LeadershipRepository
	PartyRepository      repository.PartyRepository
}

func NewMemberService(repos *repository.Repos) MemberService {
	return &memberService{
		Repository:           repos.Member,
		LeadershipRepository: repos.Leadership,
		PartyRepository:      repos.Party,
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
	return &data, nil
}

func (s *memberService) GetCareerByMemberID(memberID int) ([]map[string]interface{}, error) {
	careers, err := s.LeadershipRepository.GetCareerByMemberID(memberID)
	if err != nil {
		return nil, err
	}
	return careers, nil
}
