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

type LocalBodyService struct {
	repo repository.LocalBodyRepository
}

func NewLocalBodyService(repos *repository.Repos) *LocalBodyService {
	return &LocalBodyService{
		repo: repos.LocalBody,
	}
}

func (s *LocalBodyService) CreateLocalBody(userName string, localBodyDto *dto.LocalBodyReq) (*model.LocalBody, error) {
	localBody := &model.LocalBody{
		Name:       localBodyDto.Name,
		Code:       localBodyDto.Code,
		DistrictID: localBodyDto.DistrictID,
		Area:       localBodyDto.Area,
		BodyType:   localBodyDto.BodyType,
		TotalWards: localBodyDto.TotalWards,
		AuditFields: model.AuditFields{
			CreatedBy: userName,
		},
	}
	if exist, err := s.repo.IsExist(localBody.Name, localBody.DistrictID); err != nil {
		return nil, err
	} else if exist {
		return nil, utils.NewError("Already exists")
	}
	if err := s.repo.CreateLocalBody(localBody); err != nil {
		return nil, err
	}

	return localBody, nil
}

func (s *LocalBodyService) ListLocalBodiesByDistrict(districtId string) ([]model.LocalBody, error) {
	localBodies, err := s.repo.ListLocalBodiesByDistrict(districtId)
	if err != nil {
		return nil, err
	}
	return localBodies, nil
}

func (s *LocalBodyService) BulkCreateLocalBodies(userName string, c *gin.Context, f multipart.File) (map[string]string, map[string]string, error) {
	var localBodyReqList []dto.LocalBodyReq
	if err := json.NewDecoder(f).Decode(&localBodyReqList); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for _, localBodyReq := range localBodyReqList {
		localBody, err := s.CreateLocalBody(userName, &localBodyReq)
		if err != nil {
			createError[localBodyReq.Name] = err.Error()
		} else {
			createSuccess[localBody.Name] = "Created successfully"
		}
	}
	return createSuccess, createError, nil
}
