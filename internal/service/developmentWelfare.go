package service

import (
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

type DevelopmentWelfareService interface {
	Create(*gin.Context, dto.DevelopmentWelfareReq) (*model.DevelopmentWelfare, error)
	List() ([]model.DevelopmentWelfare, error)
	GetByID(id string) (*model.DevelopmentWelfare, error)
	GetDetailByID(string) (*dto.DevelopmentWelfareDetail, error)
}

type developmentWelfareService struct {
	Repository       repository.DevelopmentWelfareRepository
	MemberRepository repository.MemberRepository
}

func NewDevelopmentWelfareService(repos *repository.Repos) DevelopmentWelfareService {
	return &developmentWelfareService{
		Repository:       repos.DevelopmentWelfare,
		MemberRepository: repos.Member,
	}
}

func (s *developmentWelfareService) Create(ctx *gin.Context, data dto.DevelopmentWelfareReq) (*model.DevelopmentWelfare, error) {
	d := &model.DevelopmentWelfare{
		ID:                     primitive.NewObjectID(),
		Category:               data.Category,
		Type:                   data.Type,
		Name:                   data.Name,
		Description:            data.Description,
		PartyID:                data.PartyID,
		MemberID:               data.MemberID,
		StateID:                data.StateID,
		DistrictID:             data.DistrictID,
		LokSabhaID:             data.LokSabhaID,
		VidhanSabhaID:          data.VidhanSabhaID,
		Location:               data.Location,
		AnnouncementDate:       data.AnnouncementDate.Time,
		ExpectedCompletionDate: data.ExpectedCompletionDate.Time,
		ActualCompletionDate:   data.ActualCompletionDate.Time,
		Criteria:               data.Criteria,
		BeneficiariesNum:       data.BeneficiariesNum,
		Scope:                  data.Scope,
		Status:                 string(data.Status),
		Documents:              make([]model.Document, 0),
		StatusRevisions:        make([]model.StatusRevisions, 0),
		AuditFieldsBSON:        model.InsertBSONAuditFields(data.CreatedBy),
	}

	uniqeID := generateID(5)
	destinationPath := filepath.Join("developments", uniqeID)
	storage, err := storage.New(destinationPath)
	if err != nil {
		log.Errorf("Error initiating storage connection %v", err)
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			storage.DeleteAll(d.Documents)
			log.Errorf("Recovered in development welfare create from %v", r)
			panic(r)
		}
	}()
	for _, document := range data.Documents {
		fileMeta, err := storage.Save(ctx, document.File)
		if err != nil {
			log.Warnf("Error while saving file in storage %v", err)
		} else {
			d.Documents = append(d.Documents, model.Document{
				Url:      fileMeta.DownloadURL,
				FileName: fileMeta.FullPath,
				FileType: fileMeta.Extenstion,
				Remark:   document.Remark,
			})
		}
	}
	if err := s.Repository.Create(d); err != nil {
		storage.DeleteAll(d.Documents)
		return nil, err
	}
	return d, nil
}

func (s *developmentWelfareService) List() ([]model.DevelopmentWelfare, error) {
	resp, err := s.Repository.List()
	if err != nil {
		log.Errorf("Error fetching development welfare list: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *developmentWelfareService) GetByID(id string) (*model.DevelopmentWelfare, error) {
	idd, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	resp, err := s.Repository.GetPostByID(idd)
	if err != nil {
		log.Errorf("Error fetching development welfare by id: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *developmentWelfareService) GetDetailByID(id string) (*dto.DevelopmentWelfareDetail, error) {
	idd, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	resp, err := s.Repository.GetPostByID(idd)
	if err != nil {
		log.Errorf("Error fetching development welfare by id: %v", err)
		return nil, err
	}
	result := dto.DevelopmentWelfareDetail{
		DevelopmentWelfare: *resp,
	}
	var memberDto dto.Member
	memberInfo, err := s.MemberRepository.GetMemberById(resp.MemberID)
	if err != nil {
		log.Errorf("Error fetching member info by id %d", resp.MemberID)
	} else {
		if err := utils.MapToStruct(memberInfo, &memberDto); err != nil {
			log.Errorf("Error Converting member map to dto.member")
		}
	}
	result.Member = memberDto
	return &result, nil
}
