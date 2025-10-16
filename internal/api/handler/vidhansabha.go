package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type VidhanSabhaHandler struct {
	Service service.VidhanSabhaService
}

func NewVidhanSabhaHandler(serivce service.VidhanSabhaService) *VidhanSabhaHandler {
	return &VidhanSabhaHandler{
		Service: serivce,
	}
}

func (h *VidhanSabhaHandler) CreateVidhanSabha(c *gin.Context) {
	var req dto.VidhanReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create loksabha", err)
		return
	}
	utils.Success(c, "Vidhansabh created successfully", reps)
}

func (h *VidhanSabhaHandler) ListVidhanSabha(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch vidhanSabha", nil)
		return
	}
	utils.Success(c, "VidhanSabha fetched successfully", val)
}

func (h *VidhanSabhaHandler) GetByStateID(c *gin.Context) {
	stateId := c.Param("id")
	if stateId == "" {
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	stateIdInt, err := utils.StrToInt(stateId)
	if err != nil {
		utils.Error(c, 400, "Invalid state ID", nil)
		return
	}
	val, err := h.Service.GetByStateID(stateIdInt)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch vidhanSabha", nil)
		return
	}
	utils.Success(c, "VidhanSabha fetched successfully", val)
}
