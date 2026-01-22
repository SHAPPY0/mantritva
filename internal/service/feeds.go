package service

import (
	"fmt"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedService interface {
	Create(dto.FeedReq) (bool, error)
	UpdateStatus(string, dto.FeedReq) (bool, error)
	GetLatestFeed(int) (*[]dto.FeedRes, error)
	SaveReaction(dto.ReactionReq, string) (bool, error)
}

type feedService struct {
	FeedRepository          repository.FeedRepository
	NewsPressPostRepository repository.NewsPressPostRepository
	MemberRepository        repository.MemberRepository
	PartyRepository         repository.PartyRepository
	LeadershipRepository    repository.LeadershipRepository
	StateRepository         repository.StateRepository
	EntityRepository        repository.EntityRepository
	GrievanceRepository     repository.GrievanceRepository
}

func NewFeedService(repos *repository.Repos) FeedService {
	return &feedService{
		FeedRepository:          repos.Feed,
		NewsPressPostRepository: repos.NewsPressPost,
		MemberRepository:        repos.Member,
		PartyRepository:         repos.Party,
		LeadershipRepository:    repos.Leadership,
		StateRepository:         repos.State,
		EntityRepository:        repos.Entity,
		GrievanceRepository:     repos.Grievance,
	}
}

func (s *feedService) Create(feed dto.FeedReq) (bool, error) {
	f := &model.Feed{
		Type:            string(feed.Type),
		PostID:          feed.PostID,
		Status:          feed.Status,
		AuditFieldsBSON: model.InsertBSONAuditFields(),
	}
	_, err := s.FeedRepository.Create(f)
	if err != nil {
		return false, fmt.Errorf("failed to save reaction: %v", err)
	}
	return true, nil
}

func (s *feedService) GetLatestFeed(userID int) (*[]dto.FeedRes, error) {
	results := make([]dto.FeedRes, 0)
	feeds, err := s.FeedRepository.GetAll()
	if err != nil {
		log.Errorf("Error while fetching feeds: %v", err)
		return &results, err
	}
	for _, feed := range feeds {
		preparedFeed := dto.FeedRes{ID: feed.ID, Type: feed.Type}
		switch feed.Type {
		case string(utils.NEWSPRESS):
			np, err := s.GetNewsPressFeed(userID, feed.PostID)
			if err != nil {
				log.Errorf("Error while getting newsPress feed by %s: %v", feed.PostID, err)
			} else {
				preparedFeed.Feed = np
			}
		case string(utils.GRIEVANCE):
			g, err := s.GetGrievanceFeed(userID, feed.PostID)
			if err != nil {
				log.Errorf("Error while getting grievance feed by %s: %v", feed.PostID, err)
			} else {
				preparedFeed.Feed = g
			}

		default:
			log.Info("No Valid Feed Type Found")
		}
		results = append(results, preparedFeed)
	}
	return &results, nil
}

func (s *feedService) SaveReaction(req dto.ReactionReq, userIP string) (bool, error) {
	// Get existing post
	post, err := s.NewsPressPostRepository.GetPostByID(req.PostID)
	if err != nil {
		return false, fmt.Errorf("failed to get post: %v", err)
	}
	// Check for existing reaction
	reactionExist, err := s.NewsPressPostRepository.GetReactionByUserIDAndPostID(req.UserID, req.PostID)
	if err != nil {
		return false, fmt.Errorf("failed to check existing reaction: %v", err)
	}
	// Handle existing reaction
	if reactionExist != nil {
		// Same reaction type - no change needed
		if reactionExist.Reaction == string(req.Reaction) {
			return true, nil
		}

		// Different reaction - update counts and soft delete old reaction
		if err := s.NewsPressPostRepository.DeleteReaction(reactionExist.ID, req.UserID, req.PostID); err != nil {
			return false, fmt.Errorf("failed to delete old reaction: %v", err)
		}

		// Update counts based on reaction change
		if reactionExist.Reaction == string(utils.UPVOTE) {
			post.UpvoteCount--
			if req.Reaction == utils.DOWNVOTE {
				post.DownvoteCount++
			}
		} else if reactionExist.Reaction == string(utils.DOWNVOTE) {
			post.DownvoteCount--
			if req.Reaction == utils.UPVOTE {
				post.UpvoteCount++
			}
		}
	} else {
		// New reaction - just increment the appropriate count
		switch req.Reaction {
		case utils.UPVOTE:
			post.UpvoteCount++
		case utils.DOWNVOTE:
			post.DownvoteCount++
		}
	}

	// Create new reaction
	reaction := &model.PostReaction{
		PostID:          req.PostID,
		UserID:          req.UserID,
		Reaction:        string(req.Reaction),
		UserIP:          userIP,
		AuditFieldsBSON: model.InsertBSONAuditFields(),
	}
	// Save new reaction
	if err := s.NewsPressPostRepository.SaveReaction(reaction, reactionExist == nil); err != nil {
		return false, fmt.Errorf("failed to save reaction: %v", err)
	}
	// Update post vote counts
	s.updatePostVotes(post)
	return true, nil
}

func (s *feedService) updatePostVotes(newsPress *model.NewsPressPost) {
	if err := s.NewsPressPostRepository.UpdatePostVoteCountsByID(newsPress.ID, newsPress.UpvoteCount, newsPress.DownvoteCount); err != nil {
		log.Error("Error: " + err.Error())
	}
}

func (s *feedService) GetNewsPressFeed(userID int, postID primitive.ObjectID) (dto.NewsPressPostRes, error) {
	var newsPressPostRes dto.NewsPressPostRes

	np, err := s.NewsPressPostRepository.GetPostByID(postID)
	if err != nil {
		return newsPressPostRes, err
	}

	newsPressPostRes.NewsPressPost = *np
	if np.MemberId > 0 {
		member, err := s.MemberRepository.GetByID(np.MemberId)
		if err != nil {
			return newsPressPostRes, err
		}

		newsPressPostRes.Member = dto.Member{Name: member.Name, Photo: member.Photo}

		leader, err := s.LeadershipRepository.GetByMemberID(member.ID)
		if err == nil {
			if leader.Role != "" {
				newsPressPostRes.Member.Position = leader.Role
			} else {
				position := leader.Position
				newsPressPostRes.Member.Position = position
				if leader.StateID > 0 {
					state, err := s.StateRepository.GetByID(int(leader.StateID))
					if err == nil {
						newsPressPostRes.Member.Position = position + ", " + state.Name
					}
				}
			}
		}
	}
	if np.PartyId > 0 {
		party, err := s.PartyRepository.GetByID(int(np.PartyId))
		if err != nil {
			return newsPressPostRes, err
		}
		newsPressPostRes.Party = dto.Party{Name: party.Name, Abbrv: party.Abbrv, Symbol: party.Symbol}
	}
	if np.EntityId > 0 {
		ok, entity, err := s.EntityRepository.GetEntityByID(int(np.EntityId))
		if err != nil {
			return newsPressPostRes, err
		}
		if ok {
			newsPressPostRes.Entity = dto.Entity{
				Name:         entity.Name,
				Abbrv:        entity.Abbrv,
				IsCentral:    entity.IsCentral,
				StateName:    entity.StateName,
				DistrictName: entity.DistrictName,
			}
		}
	}

	if userID > 0 {
		reaction, err := s.NewsPressPostRepository.GetReactionByUserIDAndPostID(userID, np.ID)
		if err != nil {
			// fmt.Errorf("no : %v", err)
		}
		if reaction != nil {
			newsPressPostRes.Reaction = utils.PostReactionType(reaction.Reaction)
		}
	}
	return newsPressPostRes, nil
}

func (s *feedService) GetGrievanceFeed(userID int, postID primitive.ObjectID) (*dto.GrievanceRes, error) {
	var grievanceRes dto.GrievanceRes
	grievance, err := s.GrievanceRepository.GetByID(postID)
	if err != nil {
		return &grievanceRes, err
	}

	grievanceRes.Grievance = *grievance
	return &grievanceRes, nil
}

func (s *feedService) UpdateStatus(updatedBy string, feed dto.FeedReq) (bool, error) {
	ok, err := s.FeedRepository.UpdateStatus(updatedBy, feed)
	if err != nil {
		return false, err
	}
	return ok, nil
}
