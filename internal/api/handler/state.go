package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type StateHandler struct {
	Service service.StateService
}

func NewStateHandler(svc service.StateService) *StateHandler {
	return &StateHandler{
		Service: svc,
	}
}

func (h *StateHandler) CreateState(c *gin.Context) {
	var req dto.StateReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create state", err)
		return
	}
	utils.Success(c, "State created successfully", reps)
}

func (h *StateHandler) ListState(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch states", nil)
		return
	}
	utils.Success(c, "States fetched successfully", val)
}

func (h *StateHandler) GetState(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.Error(c, 400, "State name is required", nil)
		return
	}
	val, err := h.Service.Get(name)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch state "+name, err)
		return
	}
	if val == nil {
		utils.Error(c, 404, "State not found "+name, nil)
		return
	}
	utils.Success(c, "State fetched successfully", val)
}
