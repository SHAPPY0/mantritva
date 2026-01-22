package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventService interface {
	Create(dto.EventReq) (*model.Event, error)
	GetAll() ([]*model.Event, error)
}

type eventService struct {
	Repository repository.EventRepository
}

func NewEventService(repos *repository.Repos) EventService {
	return &eventService{
		Repository: repos.Event,
	}
}

func (s *eventService) Create(req dto.EventReq) (*model.Event, error) {
	e := &model.Event{
		ID:              primitive.NewObjectID(),
		Name:            req.Name,
		Description:     req.Description,
		Date:            *req.Date.Time,
		AuditFieldsBSON: model.InsertBSONAuditFields(),
	}
	if err := s.Repository.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *eventService) GetAll() ([]*model.Event, error) {
	events, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	return events, nil
}
