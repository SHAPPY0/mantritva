package service

import (
	"encoding/json"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/utils"
)

type UlbService struct {
	repo repository.UlbRepository
}

func NewULBService(repos *repository.Repos) *UlbService {
	return &UlbService{
		repo: repos.Ulb,
	}
}

func (s *UlbService) CreateULB(userName string, ulbDto *dto.ULBReq) (*model.ULB, error) {
	ulb := &model.ULB{
		Name:       ulbDto.Name,
		Code:       ulbDto.Code,
		DistrictID: ulbDto.DistrictID,
		Area:       ulbDto.Area,
		BodyType:   ulbDto.BodyType,
		TotalWards: ulbDto.TotalWards,
		AuditFields: model.AuditFields{
			CreatedBy: userName,
		},
	}
	if exist, err := s.repo.IsExist(ulb.Name, ulb.DistrictID); err != nil {
		return nil, err
	} else if exist {
		return nil, utils.NewError("Already exists")
	}
	if err := s.repo.CreateULB(ulb); err != nil {
		return nil, err
	}

	return ulb, nil
}

func (s *UlbService) ListULBsByDistrict(districtId string) ([]model.ULB, error) {
	ulbs, err := s.repo.ListULBsByDistrict(districtId)
	if err != nil {
		return nil, err
	}
	return ulbs, nil
}

func (s *UlbService) BulkCreateULB(userName string, c *gin.Context, f multipart.File) (map[string]string, map[string]string, error) {
	var ulbReqList []dto.ULBReq
	if err := json.NewDecoder(f).Decode(&ulbReqList); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for _, ulbReq := range ulbReqList {
		ulb, err := s.CreateULB(userName, &ulbReq)
		if err != nil {
			createError[ulbReq.Name] = err.Error()
		} else {
			createSuccess[ulb.Name] = "Created successfully"
		}
	}
	return createSuccess, createError, nil
}
