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

type VidhanSabhaService interface {
	List() ([]model.Vidhansabha, error)
	Create(string, dto.VidhanReq) (*model.Vidhansabha, error)
	GetByStateID(int) (*[]dto.VidhanSabhaRes, error)
	CreateBulkVidhansabha(*gin.Context, string, multipart.File) (map[string]string, map[string]string, error)
}

type vidhanSabhaService struct {
	Repository repository.VidhansabhaRepository
}

func NewVidhanSabhaService(repo repository.VidhansabhaRepository) VidhanSabhaService {
	return &vidhanSabhaService{
		Repository: repo,
	}
}

func (s *vidhanSabhaService) Create(username string, req dto.VidhanReq) (*model.Vidhansabha, error) {
	vidhanSabha := &model.Vidhansabha{
		Name:        req.Name,
		StateID:     req.StateID,
		DistrictID:  req.DistrictID,
		ACNumber:    req.ACNumber,
		AuditFields: model.AuditFields{CreatedBy: username},
	}
	if req.ACNumber > 0 {
		vidhanSabha.ACNumber = req.ACNumber
	}
	exist, err := s.Repository.GetByStateIDAndDistrictIDAndName(vidhanSabha.StateID, vidhanSabha.DistrictID, vidhanSabha.Name)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, utils.NewError("Vidhansabha already exist")
	}
	if err := s.Repository.Create(vidhanSabha); err != nil {
		return nil, err
	}
	return vidhanSabha, nil
}

func (s *vidhanSabhaService) List() ([]model.Vidhansabha, error) {
	vidhanSabhas, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return vidhanSabhas, nil
}

func (s *vidhanSabhaService) GetByStateID(stateId int) (*[]dto.VidhanSabhaRes, error) {
	vidhanSabhas, err := s.Repository.GetByStateID(stateId)
	if err != nil {
		return nil, err
	}
	return vidhanSabhas, nil
}

func (s *vidhanSabhaService) CreateBulkVidhansabha(c *gin.Context, username string, file multipart.File) (map[string]string, map[string]string, error) {
	var req []dto.VidhanReq
	if err := json.NewDecoder(file).Decode(&req); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for _, vidhansabha := range req {
		v, err := s.Create(username, vidhansabha)
		if err != nil {
			createError[vidhansabha.Name] = err.Error()
		} else if v != nil {
			createSuccess[v.Name] = "Success"
		}
	}
	return createSuccess, createError, nil
}
