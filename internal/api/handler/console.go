package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsoleHandler struct {
	Service     service.ConsoleService
	FeedService service.FeedService
}

func NewConsoleHandler(services *service.Services) *ConsoleHandler {
	return &ConsoleHandler{
		Service:     services.ConsoleService,
		FeedService: services.FeedService,
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

func (h *ConsoleHandler) GetGrievancesList(c *gin.Context) {
	resp, err := h.Service.GetGrievancesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get grievances list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetGrievanceByID(c *gin.Context) {
	grievanceID := c.Param("grievanceID")
	resp, err := h.Service.GetGrievanceByID(grievanceID)
	if err != nil {
		utils.Error(c, 500, "Failed to get grievance", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GrievanceUpdateStatus(c *gin.Context) {
	grievanceID := c.Param("grievanceID")
	var req dto.GrievanceStatus
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	resp, err := h.Service.GrievanceUpdateStatus(grievanceID, req)
	if err != nil {
		utils.Error(c, 500, "Failed to update grievance status", err)
		return
	}
	if req.Status == utils.OPEN {
		postID, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			utils.Error(c, 500, "Invalid grievance ID", err)
			return
		}
		_, err = h.FeedService.Create(dto.FeedReq{
			Type:   utils.GRIEVANCE,
			PostID: postID,
			Status: utils.Published,
		})
		if err != nil {
			utils.Error(c, 500, "Failed to create grievance feed", err)
			return
		} else {
			fmt.Printf("Grievance feed created successfully for grievance ID: %s\n", req.ID)
		}
	}
	utils.Success(c, "Grievance status updated successfully", resp)
}
