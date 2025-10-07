package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type ConsoleHandler struct {
	Service service.ConsoleService
}

func NewConsoleHandler(services *service.Services) *ConsoleHandler {
	return &ConsoleHandler{
		Service: services.ConsoleService,
	}
}

func (h *ConsoleHandler) GetUsersList(c *gin.Context) {
	resp, err := h.Service.GetUsersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get users list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetStatesList(c *gin.Context) {
	resp, err := h.Service.GetStatesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get states list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetPartiesList(c *gin.Context) {
	resp, err := h.Service.GetPartiesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get parties list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetAlliancesList(c *gin.Context) {
	resp, err := h.Service.GetAlliancesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get alliances list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetMembersList(c *gin.Context) {
	resp, err := h.Service.GetMembersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get members list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetRulersList(c *gin.Context) {
	resp, err := h.Service.GetRulersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get rulers list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetLeadershipByRulerID(c *gin.Context) {
	rulerIDStr := c.Param("rulerID")
	// Convert string to int
	rulerID, err := strconv.Atoi(rulerIDStr)
	if err != nil {
		utils.Error(c, 400, "Invalid rulerID", err)
		return
	}
	resp, err := h.Service.GetLeadershipByRulerID(rulerID)
	if err != nil {
		utils.Error(c, 500, "Failed to get rulers list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetNewsPressList(c *gin.Context) {
	resp, err := h.Service.GetNewsPressList()
	if err != nil {
		utils.Error(c, 500, "Failed to get news press list", err)
		return
	}
	utils.Success(c, "", resp)
}
