package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type RulerHandler struct {
	Service service.RulerService
}

func NewRulerHandler(serivce service.RulerService) *RulerHandler {
	return &RulerHandler{
		Service: serivce,
	}
}

func (h *RulerHandler) CreateRuler(c *gin.Context) {
	var req dto.RulerReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create ruler", err)
		return
	}
	utils.Success(c, "Ruler created successfully", reps)
}

func (h *RulerHandler) ListRulers(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch Ruler", nil)
		return
	}
	utils.Success(c, "Ruler fetched successfully", val)
}

func (h *RulerHandler) GetCentralRuler(c *gin.Context) {
	val, err := h.Service.GetCentralRuler()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch central ruler", nil)
		return
	}
	if val == nil {
		utils.Error(c, 404, "Central ruler not found", nil)
		return
	}
	utils.Success(c, "Central ruler fetched successfully", val)
}

func (h *RulerHandler) GetStateRuler(c *gin.Context) {
	stateId := c.Param("stateId")
	if stateId == "" {
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	val, err := h.Service.GetStateRuler(stateId)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch state ruler", err)
		return
	}
	if val == nil {
		utils.Error(c, 404, "State ruler not found for ID "+stateId, nil)
		return
	}
	utils.Success(c, "State ruler fetched successfully", val)
}

func (h *RulerHandler) GetCenterLeadership(c *gin.Context) {
	val, err := h.Service.GetCenterLeadership()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch center leadership", nil)
		return
	}
	if val == nil {
		utils.Error(c, 404, "Center leadership not found", nil)
		return
	}
	utils.Success(c, "Center leadership fetched successfully", val)
}
