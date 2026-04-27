package service

import (
	"encoding/json"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
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
	var ulbReq []dto.ULBReq
	if err := json.NewDecoder(f).Decode(&ulbReq); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for _, ulb := range ulbReq {
		ulb, err := s.CreateULB(userName, &ulb)
		if err != nil {
			createError[ulb.Name] = err.Error()
		} else {
			createSuccess[ulb.Name] = "Created successfully"
		}
	}
	return createSuccess, createError, nil
}
