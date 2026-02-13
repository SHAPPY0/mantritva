package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/storage"
)

type PartyService interface {
	List() ([]model.Party, error)
	Create(dto.PartyReq) (*model.Party, error)
	GetByAbbrv(abbrv string) (*model.Party, error)
	UploadSymbol(ctx *gin.Context, partyId string, file *multipart.FileHeader) (bool, error)
	AddMember(dto.AddMemberReq) (bool, error)
	GetPartyMembers(partyId string) ([]dto.PartyMemberRes, error)
	GetRandomParties() (*[]dto.Party, error)
}

type partyService struct {
	Repository       repository.PartyRepository
	MemberRepository repository.MemberRepository
}

func NewPartyService(repo *repository.Repos) PartyService {
	return &partyService{
		Repository:       repo.Party,
		MemberRepository: repo.Member,
	}
}

func (s *partyService) Create(req dto.PartyReq) (*model.Party, error) {
	party := &model.Party{
		Name:        req.Name,
		Abbrv:       req.Abbrv,
		Symbol:      req.Symbol,
		About:       req.About,
		CoreBeliefs: req.CoreBeliefs,
		Leaning:     req.Leaning,
		EcStatus:    req.EcStatus,
		Founders:    req.Founders,
		DOF:         req.DOF,
	}
	if err := s.Repository.Create(party); err != nil {
		return nil, err
	}
	return party, nil
}

func (s *partyService) List() ([]model.Party, error) {
	parties, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return parties, nil
}

func (s *partyService) GetByAbbrv(abbrv string) (*model.Party, error) {
	party, err := s.Repository.GetByAbbrv(abbrv)
	if err != nil {
		return nil, err
	}
	return party, nil
}

func (s *partyService) UploadSymbol(ctx *gin.Context, partyId string, file *multipart.FileHeader) (bool, error) {
	absoluteDestPath := filepath.Join("parties", partyId)
	storage, err := storage.New(absoluteDestPath)
	if err != nil {
		return false, err
	}
	fileMeta, err := storage.Save(ctx, file)
	if err != nil {
		return false, err
	} else {
		if err := s.Repository.UploadSymbol(partyId, fileMeta.DownloadURL); err != nil {
			_ = storage.Delete(fileMeta.AbsolutePath)
			return false, err
		}
	}
	return true, nil
}

func (s *partyService) AddMember(req dto.AddMemberReq) (bool, error) {
	partyMember := &model.PartyMembers{
		PartyID:     req.PartyId,
		MemberID:    req.MemberId,
		Position:    req.Position,
		JoiningDate: req.JoiningDate,
		ExitDate:    req.ExitDate,
	}
	if err := s.Repository.AddMember(partyMember); err != nil {
		return false, err
	}
	return true, nil
}

func (s *partyService) GetPartyMembers(partyId string) ([]dto.PartyMemberRes, error) {
	var res []dto.PartyMemberRes
	partyMembers, err := s.Repository.GetPartyMembers(partyId)
	if err != nil {
		return nil, err
	}

	for _, partyMember := range partyMembers {
		partyMemberDto := dto.PartyMemberRes{
			ID:          partyMember.ID,
			MemberID:    partyMember.MemberID,
			Position:    partyMember.Position,
			JoiningDate: partyMember.JoiningDate,
			ExitDate:    partyMember.ExitDate,
		}
		member, err := s.MemberRepository.GetByID(partyMember.MemberID)
		if err != nil {
			fmt.Printf("Failed to get member by ID %d\n", partyMember.MemberID)
		} else {
			partyMemberDto.Name = member.Name
			partyMemberDto.Photo = member.Photo
		}
		res = append(res, partyMemberDto)
	}
	//Sort By Name
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res, nil
}

func (s *partyService) GetRandomParties() (*[]dto.Party, error) {
	var results []dto.Party
	parties, err := s.Repository.GetRandomParties()
	if err != nil {
		return nil, err
	}
	for _, party := range *parties {
		results = append(results, dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol})
	}
	return &results, nil
}
