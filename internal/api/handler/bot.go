package handler

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotHandler struct {
	UserSvc               service.UserService
	FeedSvc               service.FeedService
	DevelopmentWelfareSvc service.DevelopmentWelfareService
	GrievanceSvc          service.GrievanceSerivce
	NewsPressSvc          service.NewsPressService
}

func NewBotHandler(service *service.Services) *BotHandler {
	return &BotHandler{
		UserSvc:               service.UserService,
		FeedSvc:               service.FeedService,
		DevelopmentWelfareSvc: service.DevelopmentWelfareService,
		GrievanceSvc:          service.GrievanceService,
		NewsPressSvc:          service.NewsPressService,
	}
}

func (h *BotHandler) CreateUser(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	var req dto.UserReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	user, err := h.UserSvc.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create user", err)
		return
	} else if user == nil {
		utils.Error(c, 500, "Error while creating Bot user", err)
		return
	}
	utils.Success(c, "Bot user created.", user)
}

func splitSlice(users []model.User, percent float64) ([]model.User, []model.User) {
	shuffled := make([]model.User, len(users))
	copy(shuffled, users)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Handle percent as 0-100 range (e.g., 80 = 80%)
	p := percent
	if p > 1 {
		p = p / 100
	}

	// Clamp to valid range [0, 1]
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}

	splitIdx := int(float64(len(shuffled)) * p)

	return shuffled[:splitIdx], shuffled[splitIdx:]
}

func (h *BotHandler) ReactFeeds(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}

	var req dto.FeedBotReactionReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}

	//Get List of bot users
	botUsers, err := h.UserSvc.GetBotUsers(req.TotalUsers)
	if err != nil {
		utils.Error(c, 500, "Error while fetching bot users", err)
		return
	}

	positiveUsers, negativeUsers := splitSlice(botUsers, float64(req.PositiveUsersPortion))

	feedType, feed, err := h.FeedSvc.GetPostByFeedID(req.FeedID)
	if err != nil {
		utils.Error(c, 500, "Internal Server Error", err)
		return
	}
	if feed != nil {
		var postID primitive.ObjectID
		var quesOptions []uint

		//Get questionnaires
		ques, err := h.FeedSvc.GetQuestionnairesByFeedType(feedType)
		if err != nil {
			log.Errorf("Error getting questionnaires for %s", feedType)
		} else if ques != nil {
			for _, q := range *ques {
				quesOptions = append(quesOptions, q.ID)
			}
		}
		//Get postID
		switch v := feed.(type) {
		case dto.DevelopmentWelfareShort:
			id, _ := primitive.ObjectIDFromHex(v.ID)
			postID = id
		case *dto.GrievanceShort:
			postID = v.ID
		case dto.NewsPressShort:
			postID = v.ID
		default:
			log.Errorf("Unknown type: %T\n", v)
		}
		h.createBotReaction(positiveUsers, postID, quesOptions, feedType, utils.UPVOTE)
		h.createBotReaction(negativeUsers, postID, quesOptions, feedType, utils.DOWNVOTE)
	} else {
		log.Errorf("No feed found with id %v", req.FeedID)
		utils.Error(c, 400, "No feed found with id", nil)
		return
	}

	utils.Success(c, "Bot ran successfully", nil)
}

func (h *BotHandler) createBotReaction(users []model.User, postID primitive.ObjectID, quesOptions []uint, feedType string, reaction utils.Vote) error {
	for _, user := range users {
		selectedQues, _ := pickRandom(quesOptions)
		reactionReq := dto.ReactionReq{
			UserID:          int(user.ID),
			PostID:          postID,
			Reaction:        reaction,
			FeedType:        utils.FeedType(feedType),
			SelectedOptions: selectedQues,
			UserIP:          "127.0.0.1",
		}
		h.FeedSvc.CreateReaction(reactionReq)
	}
	return nil
}

func pickRandom(items []uint) ([]int, bool) {
	if len(items) == 0 {
		return []int{1}, true
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := items[r.Intn(len(items))]
	return []int{int(i)}, true
}
