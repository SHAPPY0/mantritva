package service

import (
	"fmt"
	"math"

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
	FeedRepository               repository.FeedRepository
	NewsPressPostRepository      repository.NewsPressPostRepository
	MemberRepository             repository.MemberRepository
	PartyRepository              repository.PartyRepository
	LeadershipRepository         repository.LeadershipRepository
	StateRepository              repository.StateRepository
	EntityRepository             repository.EntityRepository
	GrievanceRepository          repository.GrievanceRepository
	DevelopmentWelfareRepository repository.DevelopmentWelfareRepository
	ReportRepository             repository.ReportRepository
}

type postSummary struct {
	ID            primitive.ObjectID
	UpvoteCount   uint
	DownvoteCount uint
}

func NewFeedService(repos *repository.Repos) FeedService {
	svc := &feedService{
		FeedRepository:               repos.Feed,
		NewsPressPostRepository:      repos.NewsPressPost,
		MemberRepository:             repos.Member,
		PartyRepository:              repos.Party,
		LeadershipRepository:         repos.Leadership,
		StateRepository:              repos.State,
		EntityRepository:             repos.Entity,
		GrievanceRepository:          repos.Grievance,
		DevelopmentWelfareRepository: repos.DevelopmentWelfare,
		ReportRepository:             repos.Report,
	}
	return svc
}

func (s *feedService) Create(feed dto.FeedReq) (bool, error) {
	f := &model.Feed{
		Type:            string(feed.Type),
		PostID:          feed.PostID,
		Status:          feed.Status,
		AuditFieldsBSON: model.InsertBSONAuditFields(""),
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
		case string(utils.DEVELOPMENT):
			dw, err := s.GetDevelopmentWelfareShort(userID, feed.PostID)
			if err != nil {
				log.Errorf("Error while getting development welfare feed by %s: %v", feed.PostID, err)
			} else {
				preparedFeed.Feed = dw
			}
		default:
			log.Info("No Valid Feed Type Found")
		}
		results = append(results, preparedFeed)
	}
	return &results, nil
}

func (s *feedService) getPostSummary(feedType utils.FeedType, postID primitive.ObjectID) (*postSummary, error) {
	switch feedType {
	case utils.NEWSPRESS:
		np, err := s.NewsPressPostRepository.GetPostByID(postID)
		if err != nil {
			return nil, err
		}
		return &postSummary{ID: np.ID, UpvoteCount: np.UpvoteCount, DownvoteCount: np.DownvoteCount}, nil
	case utils.DEVELOPMENT:
		dw, err := s.DevelopmentWelfareRepository.GetPostByID(postID)
		if err != nil {
			return nil, err
		}
		return &postSummary{ID: dw.ID, UpvoteCount: dw.UpvoteCount, DownvoteCount: dw.DownvoteCount}, nil
	default:
		return nil, fmt.Errorf("Unsupported feed type: %v", feedType)
	}
}

func (s *feedService) persistPostCounts(feedType utils.FeedType, ps *postSummary) error {
	switch feedType {
	case utils.NEWSPRESS:
		return s.NewsPressPostRepository.UpdatePostVoteCountsByID(ps.ID, ps.UpvoteCount, ps.DownvoteCount)
	case utils.DEVELOPMENT:
		return s.DevelopmentWelfareRepository.UpdatePostVoteCountsByID(ps.ID, ps.UpvoteCount, ps.DownvoteCount)
	default:
		return fmt.Errorf("Unsupported feed type: %v", feedType)
	}
}

func (s *feedService) SaveReaction(req dto.ReactionReq, userIP string) (bool, error) {
	var postID primitive.ObjectID
	var oldReactionScore int
	// Get existing post
	post, err := s.getPostSummary(req.FeedType, req.PostID)
	if err != nil {
		return false, fmt.Errorf("Failed to get post: %v", err)
	}
	postID = post.ID

	// Check for existing reaction
	reactionExist, err := s.FeedRepository.GetReactionByUserIDAndPostID(req.UserID, postID)
	if err != nil {
		return false, fmt.Errorf("Failed to check existing reaction: %v", err)
	}
	// Handle existing reaction
	if reactionExist != nil {
		// Same reaction type - no change needed
		if reactionExist.Reaction == string(req.Reaction) && utils.CompareIntSlices(reactionExist.SelectedOptions, req.SelectedOptions) {
			return true, nil
		}

		// Different reaction - update counts and soft delete old reaction
		if err := s.FeedRepository.DeleteReaction(reactionExist.ID, req.UserID, req.PostID); err != nil {
			return false, fmt.Errorf("Failed to delete old reaction: %v", err)
		}

		oldReactionScore = reactionExist.NormalizedScore

		// Update counts based on reaction change
		if reactionExist.Reaction == string(utils.UPVOTE) && req.Reaction == utils.DOWNVOTE {
			if post.UpvoteCount > 0 {
				post.UpvoteCount--
			}
			post.DownvoteCount++
		} else if reactionExist.Reaction == string(utils.DOWNVOTE) && req.Reaction == utils.UPVOTE {
			if post.DownvoteCount > 0 {
				post.DownvoteCount--
			}
			post.UpvoteCount++
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
		SelectedOptions: req.SelectedOptions,
		AuditFieldsBSON: model.InsertBSONAuditFields(""),
	}
	if req.FeedType == utils.NEWSPRESS {
		reactionScore := 0
		questionnaire, err := s.ReportRepository.GetQuestionnairesByCategory(string(utils.NEWSPRESS))
		if err != nil {
			// return err
		}
		weight := getWeight(questionnaire, *reaction)
		if reaction.Reaction == string(utils.UPVOTE) {
			reactionScore += (1 + weight)
		} else if reaction.Reaction == string(utils.DOWNVOTE) {
			reactionScore += (-1 + (weight))
		}
		MIN_POSSIBLE := -10
		MAX_POSSIBLE := 10
		//Get score in scale of 0-100
		rangeSpan := MAX_POSSIBLE - MIN_POSSIBLE
		var normalizeScore int
		if rangeSpan != 0 {
			n := ((float64(reactionScore) - float64(MIN_POSSIBLE)) / float64(rangeSpan)) * 100.0
			normalizeScore = int(math.Round(n))
		} else {
			normalizeScore = 0
		}
		reaction.PrevNormalizedScore = &oldReactionScore
		reaction.NormalizedScore = normalizeScore
	}

	// Save new reaction
	if err := s.FeedRepository.SaveReaction(reaction, reactionExist == nil); err != nil {
		return false, fmt.Errorf("Failed to save reaction: %v", err)
	}
	// Update post vote counts
	if err := s.persistPostCounts(req.FeedType, post); err != nil {
		return false, fmt.Errorf("Failed to update post counts: %v", err)
	}
	return true, nil
}

func (s *feedService) GetNewsPressFeed(userID int, postID primitive.ObjectID) (dto.NewsPressShort, error) {
	var newsPressPostRes dto.NewsPressShort

	np, err := s.NewsPressPostRepository.GetPostByID(postID)
	if err != nil {
		return newsPressPostRes, err
	}
	newsPressPostRes = dto.NewsPressShort{
		ID:            np.ID,
		Event:         np.Event,
		Statement:     np.Statement,
		StatementDate: np.StatementDate,
		Subject:       np.Subject,
		Type:          np.Type,
		Tone:          np.Tone,
		UpvoteCount:   np.UpvoteCount,
		DownvoteCount: np.DownvoteCount,
		CreatedAt:     np.CreatedAt,
	}
	fetchPartyInfo := true

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

			}
		}
		if leader != nil && leader.StateID > 0 {
			state, err := s.StateRepository.GetByID(int(leader.StateID))
			if err == nil {
				newsPressPostRes.Member.Position += ", " + state.Name
			}
		}
		if leader == nil {
			memberParty, err := s.PartyRepository.GetByMemberID(member.ID)
			if err != nil {
				log.Errorf("No memberParty found for member %d", member.ID)
			} else if memberParty != nil {
				newsPressPostRes.Member.Position = memberParty.Position
				newsPressPostRes.Party = dto.Party{Name: memberParty.Name, Abbrv: memberParty.Abbrv, Symbol: memberParty.Symbol}
				fetchPartyInfo = false
			}
		}
	}
	if fetchPartyInfo && np.PartyId > 0 {
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
		reaction, _ := s.FeedRepository.GetReactionByUserIDAndPostID(userID, np.ID)
		if reaction != nil {
			newsPressPostRes.Reaction = utils.Vote(reaction.Reaction)
		}
	}
	//Get Questionnaires for this feed
	if ques, _ := s.ReportRepository.GetQuestionnairesByCategory(string(utils.NEWSPRESS)); ques != nil {
		newsPressPostRes.Questionnaire = new([]dto.Questionnaire)
		for _, q := range *ques {
			*newsPressPostRes.Questionnaire = append(*newsPressPostRes.Questionnaire, dto.Questionnaire{
				ID:       q.ID,
				Title:    q.Title,
				VoteType: &q.VoteType,
			})
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

func (s *feedService) GetDevelopmentWelfareShort(userID int, postID primitive.ObjectID) (*dto.DevelopmentWelfareShort, error) {
	var devWelRes dto.DevelopmentWelfareShort
	devWel, err := s.DevelopmentWelfareRepository.GetPostByID(postID)
	if err != nil {
		return &devWelRes, err
	}
	devWelFeed := dto.DevelopmentWelfareShort{
		ID:               devWel.ID.Hex(),
		Category:         devWel.Category,
		Type:             devWel.Type,
		Name:             devWel.Name,
		Description:      devWel.Description,
		AnnouncementDate: utils.DateOnly{Time: &devWel.AnnouncementDate},
		UpvoteCount:      devWel.UpvoteCount,
		DownvoteCount:    devWel.DownvoteCount,
		Status:           utils.DevelopmentStatus(devWel.Status),
		CreatedAt:        utils.DateOnly{Time: &devWel.CreatedAt},
	}
	member, err := s.MemberRepository.GetByID(devWel.MemberID)
	if err == nil {
		devWelFeed.Member = dto.Member{
			Name:  member.Name,
			Photo: member.Photo,
		}
		leader, err := s.LeadershipRepository.GetByMemberID(member.ID)
		if err == nil {
			if leader.Role != "" {
				devWelFeed.Member.Position = leader.Role
			} else {
				position := leader.Position
				devWelFeed.Member.Position = position
				if leader.StateID > 0 {
					state, err := s.StateRepository.GetByID(int(leader.StateID))
					if err == nil {
						devWelFeed.Member.Position = position + ", " + state.Name
					}
				}
			}
		}
	} else {
		log.Errorf("Error fetching member for development welfare feed: %v", err)
	}

	party, err := s.PartyRepository.GetByID(int(devWel.PartyID))
	if err == nil {
		devWelFeed.Party = dto.Party{
			Name:   party.Name,
			Abbrv:  party.Abbrv,
			Symbol: party.Symbol,
		}
	} else {
		log.Errorf("Error fetching party for development welfare feed: %v", err)
	}

	// Fetch reaction if userID is provided
	if userID > 0 {
		reaction, _ := s.FeedRepository.GetReactionByUserIDAndPostID(userID, postID)
		if reaction != nil {
			devWelFeed.Reaction = utils.Vote(reaction.Reaction)
		}
	}
	devWelRes = devWelFeed
	return &devWelRes, nil
}

func (s *feedService) UpdateStatus(updatedBy string, feed dto.FeedReq) (bool, error) {
	ok, err := s.FeedRepository.UpdateStatus(updatedBy, feed)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func getWeight(questionnaires *[]model.Questionnaire, reaction model.PostReaction) int {
	weight := 0
	for _, q := range *questionnaires {
		if q.ID == uint(reaction.SelectedOptions[0]) {
			weight = q.Weight
		}
	}
	return weight
}
