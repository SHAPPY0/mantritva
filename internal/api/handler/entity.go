package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type EntityHandler struct {
	Service service.EntityService
}

func NewEntityHandler(serivce service.EntityService) *EntityHandler {
	return &EntityHandler{
		Service: serivce,
	}
}

func (h *EntityHandler) CreateEntity(c *gin.Context) {
	var req dto.EntityReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create entity", err)
		return
	}
	utils.Success(c, "Entity created successfully", reps)
}

func (h *EntityHandler) ListEntity(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch entity", nil)
		return
	}
	utils.Success(c, "Entity fetched successfully", val)
}
