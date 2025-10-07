package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type CountryHandler struct {
	Service service.CountryService
}

func NewCountryHandler(serivce service.CountryService) *CountryHandler {
	return &CountryHandler{
		Service: serivce,
	}
}

func (h *CountryHandler) CreateCountry(c *gin.Context) {
	var req dto.CountryReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create country", err)
		return
	}
	utils.Success(c, "Country created successfully", reps)
}

func (h *CountryHandler) ListCountry(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch countries", nil)
	}
	utils.Success(c, "Countries fetched successfully", val)
}
