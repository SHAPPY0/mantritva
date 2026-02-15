package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsPressService interface {
	Create(dto.NewsPressReq) (*model.NewsPressPost, error)
	GetPostsByPartyID(dto.Pagination, int) (*[]dto.NewsPressPostRes, error)
	GetDetailByID(int, string) (*dto.NewsPressPostRes, error)
}

type newsPressService struct {
	Repository     repository.NewsPressPostRepository
	MemberRepo     repository.MemberRepository
	PartyRepo      repository.PartyRepository
	LeadershipRepo repository.LeadershipRepository
	StateRepo      repository.StateRepository
	EntityRepo     repository.EntityRepository
	FeedRepo       repository.FeedRepository
}

func NewNewsPressService(repos *repository.Repos) NewsPressService {
	return &newsPressService{
		Repository:     repos.NewsPressPost,
		MemberRepo:     repos.Member,
		PartyRepo:      repos.Party,
		LeadershipRepo: repos.Leadership,
		StateRepo:      repos.State,
		EntityRepo:     repos.Entity,
		FeedRepo:       repos.Feed,
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

func (s *newsPressService) GetDetailByID(userID int, id string) (*dto.NewsPressPostRes, error) {
	postId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var result dto.NewsPressPostRes
	post, err := s.Repository.GetPostByID(postId)
	if err != nil {
		return nil, err
	}
	if post.MemberId > 0 {
		memberInfo, err := s.MemberRepo.GetMemberById(post.MemberId)
		if err != nil {
			log.Errorf("Error fetching member info by id %d", post.MemberId)
			return nil, err
		} else {
			if err := utils.MapToStruct(memberInfo, &result.Member); err != nil {
				log.Errorf("Error converting member map to dto.member")
			}
			if !utils.IsEmpty(result.Member.Party) {
				result.Party = dto.Party{
					Name:   result.Member.Party.Name,
					Abbrv:  result.Member.Party.Abbrv,
					Symbol: result.Member.Party.Symbol,
				}
			}
		}
	}
	if post.EntityId > 0 {
		ok, entity, err := s.EntityRepo.GetEntityByID(int(post.EntityId))
		if err != nil {
			log.Errorf("Error getting entity by id %d", post.EntityId)
		}
		if ok {
			result.Entity = dto.Entity{
				Name:         entity.Name,
				Abbrv:        entity.Abbrv,
				IsCentral:    entity.IsCentral,
				StateName:    entity.StateName,
				DistrictName: entity.DistrictName,
			}
		}
	}
	if userID > 0 {
		reaction, _ := s.FeedRepo.GetReactionByUserIDAndPostID(userID, post.ID)
		if reaction != nil {
			result.Reaction = utils.Vote(reaction.Reaction)
		}
	}
	result.NewsPressPost = *post
	return &result, nil

}
