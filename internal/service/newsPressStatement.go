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
	Repository repository.NewsPressPostRepository
	MemberRepo repository.MemberRepository
	PartyRepo  repository.PartyRepository
}

func NewNewsPressService(repos *repository.Repos) NewsPressService {
	return &newsPressService{
		Repository: repos.NewsPressPost,
		MemberRepo: repos.Member,
		PartyRepo:  repos.Party,
	}
}

func (s *newsPressService) Create(req dto.NewsPressReq) (*model.NewsPressPost, error) {
	np := &model.NewsPressPost{
		ID:              primitive.NewObjectID(),
		MemberId:        req.MemberID,
		PartyId:         req.PartyID,
		Subject:         string(req.Subject),
		Type:            string(req.Type),
		Tone:            string(req.Tone),
		Statement:       req.Statement,
		StatementDate:   *req.StatementDate.Time,
		Status:          string(req.Status),
		AuditFieldsBSON: model.InsertBSONAuditFields(),
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
