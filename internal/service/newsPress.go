package service

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsPressService interface {
	Create(dto.NewsPressReq) (*model.NewsPressPost, error)
	CreateWithMedia(*gin.Context, dto.NewsPressReq, auth.AuthUser) (*model.NewsPressPost, error)
	GetPostsByPartyID(dto.Pagination, int) (*[]dto.NewsPressPostRes, error)
	GetDetailByID(int, string) (*dto.NewsPressPostRes, error)
	GetByUserID(int, dto.Pagination) (*[]dto.NewsPressShort, error)
}

type newsPressService struct {
	Repository     repository.NewsPressPostRepository
	MemberRepo     repository.MemberRepository
	PartyRepo      repository.PartyRepository
	LeadershipRepo repository.LeadershipRepository
	StateRepo      repository.StateRepository
	EntityRepo     repository.EntityRepository
	FeedRepo       repository.FeedRepository
	ReactionRepo   repository.ReactionRepository
	ReportRepo     repository.ReportRepository
	EventRepo      repository.EventRepository
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
		ReactionRepo:   repos.Reaction,
		ReportRepo:     repos.Report,
		EventRepo:      repos.Event,
	}
}

func (s *newsPressService) Create(req dto.NewsPressReq) (*model.NewsPressPost, error) {
	np := &model.NewsPressPost{
		ID:              primitive.NewObjectID(),
		Event:           &req.Event,
		MemberID:        *req.MemberID,
		PartyId:         *req.PartyID,
		Subject:         string(*req.Subject),
		Type:            string(*req.Type),
		Tone:            string(*req.Tone),
		Statement:       req.Statement,
		StatementDate:   *req.StatementDate.Time,
		Status:          string(req.Status),
		Source:          req.Source,
		AuditFieldsBSON: model.InsertBSONAuditFields(""),
	}
	if err := s.Repository.Create(np); err != nil {
		return nil, err
	}
	return np, nil
}

func (s *newsPressService) CreateWithMedia(ctx *gin.Context, req dto.NewsPressReq, user auth.AuthUser) (*model.NewsPressPost, error) {
	np := &model.NewsPressPost{
		ID:              primitive.NewObjectID(),
		Event:           &req.Event,
		MemberID:        *req.MemberID,
		PartyId:         *req.PartyID,
		EntityId:        *req.EntityID,
		Subject:         string(*req.Subject),
		Type:            string(*req.Type),
		Tone:            string(*req.Tone),
		Statement:       req.Statement,
		StatementDate:   *req.StatementDate.Time,
		Status:          string(req.Status),
		Source:          req.Source,
		AuditFieldsBSON: model.InsertBSONAuditFields(fmt.Sprintf("%d", user.ID)),
	}
	destinationPath := filepath.Join("newspress", np.ID.Hex())
	storage, err := storage.New(destinationPath)
	if err != nil {
		log.Infof("Error initiating storage connection %v", err)
		return nil, err
	}

	for i, mediaReq := range req.Media {
		fileMeta, err := storage.Save(ctx, mediaReq)
		if err != nil {
			log.Warnf("Error while storing file %v", err)
		} else {
			media := model.NewsPressMedia{
				URL:  fileMeta.DownloadURL,
				Type: fileMeta.Type,
			}
			if i == 0 {
				media.IsMain = true
			}
			np.Media = append(np.Media, media)
		}
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
		if post.MemberID > 0 {
			member, err := s.MemberRepo.GetByID(post.MemberID)
			if err == nil && member != nil {
				newsPresPost.Member = dto.Member{Name: member.Name, Photo: member.Photo}
				leader, err := s.LeadershipRepo.GetByMemberID(member.ID)
				if err == nil {
					if leader.Role != "" {
						newsPresPost.Member.Position = leader.Role
					} else {
						position := leader.Position
						newsPresPost.Member.Position = position.Label()
						if leader.StateID > 0 {
							state, err := s.StateRepo.GetByID(int(leader.StateID))
							if err == nil {
								newsPresPost.Member.Position = position.Label() + ", " + state.Name
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
	post, err := s.Repository.GetByID(postId)
	if err != nil {
		return nil, err
	}
	if post.MemberID > 0 {
		memberInfo, err := s.MemberRepo.GetMemberById(post.MemberID)
		if err != nil {
			log.Errorf("Error fetching member info by id %d", post.MemberID)
			return nil, err
		} else {
			if err := utils.MapToStruct(memberInfo, &result.Member); err != nil {
				log.Errorf("Error converting member map to dto.member")
			} else {
				utils.TransformPositionLabel(&result.Member.Position)
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
		ok, entity, err := s.EntityRepo.GetByID(int(post.EntityId))
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
		reaction, _ := s.ReactionRepo.GetByUserIDAndPostID(userID, post.ID)
		if reaction != nil {
			result.Reaction = utils.Vote(reaction.Reaction)
			if reaction.LockedAt != nil && time.Now().After(*reaction.LockedAt) {
				result.ReactionLocked = true
			}
		}
	}

	//Get questionnair
	if ques, _ := s.ReportRepo.GetQuestionnairesByCategory(string(utils.NEWSPRESS)); ques != nil {
		result.Questionnaire = new([]dto.Questionnaire)
		for _, q := range *ques {
			*result.Questionnaire = append(*result.Questionnaire, dto.Questionnaire{
				ID:       q.ID,
				Title:    q.Title,
				VoteType: &q.VoteType,
			})
		}
	}
	result.NewsPressPost = *post
	return &result, nil

}

func (s *newsPressService) GetByUserID(userID int, pagination dto.Pagination) (*[]dto.NewsPressShort, error) {
	var newsPressPostRes []dto.NewsPressShort
	posts, err := s.Repository.GetByUserID(userID, pagination)
	if err != nil {
		return nil, err
	}
	for _, post := range *posts {
		var newsPresPost = dto.NewsPressShort{
			ID:            post.ID,
			Statement:     post.Statement,
			StatementDate: post.StatementDate,
			Subject:       post.Subject,
			Type:          post.Type,
			Tone:          post.Tone,
			UpvoteCount:   post.UpvoteCount,
			DownvoteCount: post.DownvoteCount,
			CreatedAt:     post.CreatedAt,
			Status:        utils.PostStatus(post.Status),
			Media:         post.Media,
		}
		if post.MemberID > 0 {
			member, err := s.MemberRepo.GetByID(post.MemberID)
			if err == nil && member != nil {
				newsPresPost.Member = dto.Member{Name: member.Name, Photo: member.Photo}
				leader, err := s.LeadershipRepo.GetByMemberID(member.ID)
				if err == nil {
					if leader.Role != "" {
						newsPresPost.Member.Position = leader.Role
					} else {
						position := leader.Position
						newsPresPost.Member.Position = position.Label()
						if leader.StateID > 0 {
							state, err := s.StateRepo.GetByID(int(leader.StateID))
							if err == nil {
								newsPresPost.Member.Position = position.Label() + ", " + state.Name
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
		if post.EntityId > 0 {
			ok, entity, err := s.EntityRepo.GetByID(int(post.EntityId))
			if err != nil {
				log.Errorf("Error getting entity by id %d", post.EntityId)
			}
			if ok {
				newsPresPost.Entity = dto.Entity{
					Name:         entity.Name,
					Abbrv:        entity.Abbrv,
					IsCentral:    entity.IsCentral,
					StateName:    entity.StateName,
					DistrictName: entity.DistrictName,
				}
			}
		}
		if post.Event != nil {
			event, err := s.EventRepo.GetByID(*post.Event)
			if err != nil {
				log.Errorf("Error getting event by id %s", *post.Event)
			}
			if event != nil {
				newsPresPost.Event = &event.Name
			}
		}
		newsPressPostRes = append(newsPressPostRes, newsPresPost)
	}
	return &newsPressPostRes, nil
}
