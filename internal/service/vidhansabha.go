package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type VidhanSabhaService interface {
	List() ([]model.Vidhansabha, error)
	Create(dto.VidhanReq) (*model.Vidhansabha, error)
}

type vidhanSabhaService struct {
	Repository repository.VidhansabhaRepository
}

func NewVidhanSabhaService(repo repository.VidhansabhaRepository) VidhanSabhaService {
	return &vidhanSabhaService{
		Repository: repo,
	}
}

func (s *vidhanSabhaService) Create(req dto.VidhanReq) (*model.Vidhansabha, error) {
	vidhanSabha := &model.Vidhansabha{
		Name:       req.Name,
		LoksabhaID: req.LoksabhaID,
		ACNumber:   req.ACNumber,
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
