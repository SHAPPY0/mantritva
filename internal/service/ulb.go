package service

import (
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
		ShortName:  ulbDto.ShortName,
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
