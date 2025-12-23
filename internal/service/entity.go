package service

import (
	"log"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

type EntityService interface {
	List() ([]dto.EntityRes, error)
	Create(dto.EntityReq) (*model.Entity, error)
	GetEntityByName(string) (*dto.EntityRes, error)
	GetPostsByEntityID(dto.Pagination, int) (*[]dto.NewsPressPostRes, error)
}

type entityService struct {
	Repository          repository.EntityRepository
	NewsPressRepository repository.NewsPressPostRepository
}

func NewEntityService(repos *repository.Repos) EntityService {
	return &entityService{
		Repository:          repos.Entity,
		NewsPressRepository: repos.NewsPressPost,
	}
}

func (s *entityService) Create(req dto.EntityReq) (*model.Entity, error) {
	entity := &model.Entity{
		Name:  req.Name,
		Abbrv: req.Abbrv,
		// IsCentral:  *req.IsCentral,
		StateID:    req.StateID,
		DistrictID: req.DistrictID,
	}
	if err := s.Repository.Create(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *entityService) List() ([]dto.EntityRes, error) {
	entities, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (s *entityService) GetEntityByName(name string) (*dto.EntityRes, error) {
	ok, entity, err := s.Repository.GetEntityByName(name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return entity, nil
}

func (s *entityService) GetPostsByEntityID(pagination dto.Pagination, entityID int) (*[]dto.NewsPressPostRes, error) {
	var newsPressPostRes []dto.NewsPressPostRes
	posts, err := s.NewsPressRepository.GetPostsByEntityID(pagination, entityID)
	if err != nil {
		return nil, err
	}
	for _, post := range *posts {
		dtoPost := dto.NewsPressPostRes{NewsPressPost: post}
		ok, entity, err := s.Repository.GetEntityByID(int(post.EntityId))
		if err != nil {
			log.Println("Error fetching entity: ", err)
		}
		if ok {
			dtoEntity := dto.Entity{
				Name:         entity.Name,
				Abbrv:        entity.Abbrv,
				IsCentral:    entity.IsCentral,
				StateName:    entity.StateName,
				DistrictName: entity.DistrictName,
			}
			dtoPost.Entity = dtoEntity
		}
		newsPressPostRes = append(newsPressPostRes, dtoPost)
	}
	return &newsPressPostRes, nil
}
