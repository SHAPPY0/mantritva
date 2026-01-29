package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type NewsPressHandler struct {
	Service     service.NewsPressService
	FeedService service.FeedService
}

func NewNewsPressHandler(services *service.Services) *NewsPressHandler {
	return &NewsPressHandler{
		Service:     services.NewsPressService,
		FeedService: services.FeedService,
	}
}

func (h *NewsPressHandler) CreateNewsPress(c *gin.Context) {
	var req dto.NewsPressReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	resp, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create newsPress", err)
		return
	}
	feed := dto.FeedReq{
		Type:   utils.NEWSPRESS,
		PostID: resp.ID,
		Status: utils.PostStatus(resp.Status),
	}
	ok, err := h.FeedService.Create(feed)
	if err != nil {
		utils.Error(c, 500, "Failed to create newsPress feed record", err)
		return
	}
	if ok {
		utils.Success(c, "NewsPress created successfully", resp)
	} else {
		utils.Error(c, 500, "NewsPress created but error while feed creation", err)
	}

}

func (h *NewsPressHandler) GetPostsByPartyID(c *gin.Context) {
	pagination := dto.NewPagination(c)
	partyID := c.Param("partyID")
	partyId, err := utils.StrToInt(partyID)
	if err != nil {
		utils.Error(c, 400, "Invalid party ID", err)
		return
	}

	posts, err := h.Service.GetPostsByPartyID(pagination, partyId)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch newsPress posts by party ID", err)
		return
	}
	utils.Success(c, "NewsPress posts fetched successfully", posts)
}
