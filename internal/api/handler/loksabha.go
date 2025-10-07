package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LoksabhaHandler struct {
	Service service.LoksabhaService
}

func NewLoksabhaHandler(serivce service.LoksabhaService) *LoksabhaHandler {
	return &LoksabhaHandler{
		Service: serivce,
	}
}

func (h *LoksabhaHandler) CreateLoksabha(c *gin.Context) {
	var req dto.LoksabhaReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create loksabha", err)
		return
	}
	utils.Success(c, "Loksabha created successfully", reps)
}

func (h *LoksabhaHandler) ListLoksabha(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch loksabha", nil)
		return
	}
	utils.Success(c, "Loksabha fetched successfully", val)
}
