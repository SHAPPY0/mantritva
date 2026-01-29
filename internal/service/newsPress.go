package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsPressService interface {
	Create(dto.NewsPressReq) (*model.NewsPressPost, error)
	GetPostsByPartyID(dto.Pagination, int) (*[]dto.NewsPressPostRes, error)
}

type newsPressService struct {
	Repository     repository.NewsPressPostRepository
	MemberRepo     repository.MemberRepository
	PartyRepo      repository.PartyRepository
	LeadershipRepo repository.LeadershipRepository
	StateRepo      repository.StateRepository
}

func NewNewsPressService(repos *repository.Repos) NewsPressService {
	return &newsPressService{
		Repository:     repos.NewsPressPost,
		MemberRepo:     repos.Member,
		PartyRepo:      repos.Party,
		LeadershipRepo: repos.Leadership,
		StateRepo:      repos.State,
	}
}

func (s *newsPressService) Create(req dto.NewsPressReq) (*model.NewsPressPost, error) {
	np := &model.NewsPressPost{
		ID:              primitive.NewObjectID(),
		Event:           &req.Event,
		MemberId:        req.MemberID,
		PartyId:         req.PartyID,
		Subject:         string(req.Subject),
		Type:            string(req.Type),
		Tone:            string(req.Tone),
		Statement:       req.Statement,
		StatementDate:   *req.StatementDate.Time,
		Status:          string(req.Status),
		AuditFieldsBSON: model.InsertBSONAuditFields(""),
	}
	if err := s.Repository.Create(np); err != nil {
		return nil, err
	}
	return np, nil
}

func (s *newsPressService) GetPostsByPartyID(pagination dto.Pagination, partyID int) (*[]dto.NewsPressPostRes, error) {
	var newsPressPostRes []dto.NewsPressPostRes
	posts, err := s.Repository.GetPostsByPartyID(pagination, partyID)
	if err != nil {
		return nil, err
	}
	for _, post := range *posts {
		var newsPresPost dto.NewsPressPostRes
		newsPresPost.NewsPressPost = post
		if post.MemberId > 0 {
			member, err := s.MemberRepo.GetByID(post.MemberId)
			if err == nil && member != nil {
				newsPresPost.Member = dto.Member{Name: member.Name, Photo: member.Photo}
				leader, err := s.LeadershipRepo.GetByMemberID(member.ID)
				if err == nil {
					if leader.Role != "" {
						newsPresPost.Member.Position = leader.Role
					} else {
						position := leader.Position
						newsPresPost.Member.Position = position
						if leader.StateID > 0 {
							state, err := s.StateRepo.GetByID(int(leader.StateID))
							if err == nil {
								newsPresPost.Member.Position = position + ", " + state.Name
							}
						}
					}
				}
			}
		}
		if post.PartyId > 0 {
			party, err := s.PartyRepo.GetByID(int(post.PartyId))
			if err == nil && party != nil {
				newsPresPost.Party = dto.Party{Name: party.Name, Symbol: party.Symbol, Abbrv: party.Abbrv}
			}
		}
		newsPressPostRes = append(newsPressPostRes, newsPresPost)
	}
	return &newsPressPostRes, nil
}
