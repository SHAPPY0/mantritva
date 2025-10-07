package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type DistrictHandler struct {
	Service service.DistrictService
}

func NewDistrictHandler(serivce service.DistrictService) *DistrictHandler {
	return &DistrictHandler{
		Service: serivce,
	}
}

func (h *DistrictHandler) CreateDistrict(c *gin.Context) {
	var req dto.DistrictReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create district", err)
		return
	}
	utils.Success(c, "District created successfully", reps)
}

func (h *DistrictHandler) CreateBulkDistrict(c *gin.Context) {
	var req []dto.DistrictReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	for i := 0; i < len(req); i++ {
		_, err := h.Service.Create(req[i])
		if err != nil {
			fmt.Printf("Error creating district %s: %v\n", req[i].Name, err)
		}
	}
	utils.Success(c, "Districts created successfully", req)
}

func (h *DistrictHandler) ListDistrict(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch district", nil)
		return
	}
	utils.Success(c, "District fetched successfully", val)
}

func (h *DistrictHandler) GetDistrictByState(c *gin.Context) {
	stateId := c.Param("stateId")
	if stateId == "" {
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	val, err := h.Service.GetByState(stateId)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch district for state "+stateId, err)
		return
	}
	if val == nil {
		utils.Error(c, 404, "District not found for state "+stateId, nil)
		return
	}
	utils.Success(c, "District fetched successfully", val)
}
