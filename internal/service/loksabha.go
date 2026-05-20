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

type LoksabhaService interface {
	List() ([]model.Loksabha, error)
	Create(string, dto.LoksabhaReq) (*model.Loksabha, error)
	GetByStateID(stateId int) ([]model.Loksabha, error)
	CreateBulkLoksabha(*gin.Context, string, multipart.File) (map[string]string, map[string]string, error)
}

type loksabhaService struct {
	Repository repository.LoksabhaRepository
}

func NewLoksabhaService(repo repository.LoksabhaRepository) LoksabhaService {
	return &loksabhaService{
		Repository: repo,
	}
}

func (s *loksabhaService) Create(username string, req dto.LoksabhaReq) (*model.Loksabha, error) {
	loksabha := &model.Loksabha{
		Name:       req.Name,
		StateID:    req.StateID,
		DistrictID: req.DistrictID,
		PCNumber:   req.PCNumber,
		AuditFields: model.AuditFields{
			CreatedBy: username,
		},
	}
	existLoksabha, err := s.Repository.GetByStateIDAndName(loksabha.StateID, loksabha.Name)
	if err != nil {
		return nil, err
	}
	if existLoksabha != nil {
		return nil, utils.NewError("Loksabha already exist")
	}
	if err := s.Repository.Create(loksabha); err != nil {
		return nil, err
	}
	return loksabha, nil
}

func (s *loksabhaService) List() ([]model.Loksabha, error) {
	loksabhas, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return loksabhas, nil
}

func (s *loksabhaService) GetByStateID(stateId int) ([]model.Loksabha, error) {
	loksabhas, err := s.Repository.GetByStateID(stateId)
	if err != nil {
		return nil, err
	}
	return loksabhas, nil
}

func (s *loksabhaService) CreateBulkLoksabha(c *gin.Context, username string, file multipart.File) (map[string]string, map[string]string, error) {
	var req []dto.LoksabhaReq
	if err := json.NewDecoder(file).Decode(&req); err != nil {
		panic(err)
	}
	createSuccess := map[string]string{}
	createError := map[string]string{}
	for _, loksabha := range req {
		l, err := s.Create(username, loksabha)
		if err != nil {
			createError[loksabha.Name] = err.Error()
		} else if l != nil {
			createSuccess[l.Name] = "Success"
		}
	}
	return createSuccess, createError, nil
}
