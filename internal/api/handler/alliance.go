package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type AllianceHandler struct {
	Service service.AllianceService
}

func NewAllianceHandler(serivce service.AllianceService) *AllianceHandler {
	return &AllianceHandler{
		Service: serivce,
	}
}

func (h *AllianceHandler) CreateAlliance(c *gin.Context) {
	var req dto.AllianceReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create alliance", err)
		return
	}
	utils.Success(c, "Alliance created successfully", reps)
}

func (h *AllianceHandler) ListAlliance(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch alliance", nil)
		return
	}
	utils.Success(c, "Alliance fetched successfully", val)
}

func (h *AllianceHandler) AddPartyToAlliance(c *gin.Context) {
	var req dto.AlliancePartyReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	err := h.Service.AddPartyToAlliance(req)
	if err != nil {
		utils.Error(c, 500, "Failed to add party to alliance", err)
		return
	}
	utils.Success(c, "Party added to alliance successfully", nil)
}

func (h *AllianceHandler) GetAllianceByAbbrv(c *gin.Context) {
	abbrv := c.Param("alliance")
	if abbrv == "" {
		utils.Error(c, 400, "Alliance abbreviation is required", nil)
		return
	}
	alliance, err := h.Service.GetAllianceByAbbrv(abbrv)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch alliance by abbreviation", err)
		return
	}
	utils.Success(c, "Alliance fetched successfully", alliance)
}
