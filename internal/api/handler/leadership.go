package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LeadershipHandler struct {
	Service service.LeadershipService
}

func NewLeadershipHandler(service service.LeadershipService) *LeadershipHandler {
	return &LeadershipHandler{
		Service: service,
	}
}

func (h *LeadershipHandler) Create(c *gin.Context) {
	var req dto.LeadershipReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	resp, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create leadership", err)
		return
	}
	utils.Success(c, "Leadership created successfully", resp)
}
