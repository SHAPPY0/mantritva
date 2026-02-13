package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type FeedHandler struct {
	Service                  service.FeedService
	ReportService            service.ReportService
	FeedReactionComputeQueue chan dto.ReactionReq
}

func NewFeedHandler(service *service.Services) *FeedHandler {
	h := &FeedHandler{
		Service:                  service.FeedService,
		ReportService:            service.ReportService,
		FeedReactionComputeQueue: make(chan dto.ReactionReq, 100),
	}
	return h
}

func (h *FeedHandler) GetLatestFeed(c *gin.Context) {
	userID, err := utils.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	feed, err := h.Service.GetLatestFeed(userID)
	if err != nil {
		utils.Error(c, 500, "Failed to load feed", err)
		return
	}
	utils.Success(c, "", feed)
}

func (h *FeedHandler) SaveReaction(c *gin.Context) {
	userID, err := utils.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	var req dto.ReactionReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req.UserID = userID
	resp, err := h.Service.SaveReaction(req, c.ClientIP())
	if err != nil {
		utils.Error(c, 500, "Failed to react post", err)
		return
	}
	if resp {
		utils.Success(c, "", resp)
	} else {
		utils.Error(c, 500, "Failed to vote post", nil)
	}

}
