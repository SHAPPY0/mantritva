package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type EventHandler struct {
	Service service.EventService
}

func NewEventHandler(service *service.Services) *EventHandler {
	return &EventHandler{
		Service: service.EventService,
	}
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.EventReq
	_, err := utils.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 500, "Invalid request body", err)
		return
	}

	event, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create event", err)
		return
	}

	utils.Success(c, "", event)
}

func (h *EventHandler) GetAll(c *gin.Context) {
	_, err := utils.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	events, err := h.Service.GetAll()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch events", err)
		return
	}

	utils.Success(c, "", events)
}
