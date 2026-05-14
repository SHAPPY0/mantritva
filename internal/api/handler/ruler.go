package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/log"
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
		log.Errorf("Invalid request data for create ruler %+v", err)
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		log.Errorf("Failed to create ruler %+v", err)
		utils.Error(c, 500, "Failed to create ruler", err)
		return
	}
	utils.Success(c, "Ruler created successfully", reps)
}

func (h *RulerHandler) ListRulers(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		log.Errorf("Failed to fetch ruler %+v", err)
		utils.Error(c, 500, "Failed to fetch Ruler", nil)
		return
	}
	utils.Success(c, "Ruler fetched successfully", val)
}

func (h *RulerHandler) GetCentralRuler(c *gin.Context) {
	val, err := h.Service.GetCentralRuler()
	if err != nil {
		log.Errorf("Failed to fetch central ruler %+v", err)
		utils.Error(c, 500, "Failed to fetch central ruler", nil)
		return
	}
	if val == nil {
		log.Errorf("Central ruler not found")
		utils.Error(c, 500, "Failed to fetch central ruler", nil)
		return
	}
	utils.Success(c, "Central ruler fetched successfully", val)
}

func (h *RulerHandler) GetCenterLeadership(c *gin.Context) {
	val, err := h.Service.GetCenterLeadership()
	if err != nil {
		log.Errorf("Failed to fetch center leadership %+v", err)
		utils.Error(c, 500, "Failed to fetch center leadership", nil)
		return
	}
	utils.Success(c, "Center leadership fetched successfully", val)
}

func (h *RulerHandler) GetStateRuler(c *gin.Context) {
	stateId := c.Param("stateId")
	if stateId == "" {
		log.Errorf("StateID is required to fetch stateRuler by stateID")
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	val, err := h.Service.GetStateRuler(stateId)
	if err != nil {
		log.Errorf("Failed to fetch state ruler %+v", err)
		utils.Error(c, 500, "Failed to fetch state ruler", err)
		return
	}
	utils.Success(c, "State ruler fetched successfully", val)
}

func (h *RulerHandler) GetRajyasabhaByStateID(c *gin.Context) {
	stateId := c.Param("stateId")
	if stateId == "" {
		log.Errorf("StateID is required to fetch rajyasabha by stateID")
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	var results []dto.LeadershipRes
	val, err := h.Service.GetRajyasabhaByStateID(stateId)
	if err != nil {
		log.Errorf("Failed to fetch rajyasabha members %+v", err)
		utils.Error(c, 500, "Failed to fetch rajyasabha members", err)
		return
	}
	for _, leader := range val {
		leadership := dto.LeadershipRes{}
		if err := utils.MapToStruct(leader, &leadership); err != nil {
			log.Errorf("Failed to parse rs member: %+v", err)
		} else {
			utils.TransformPositionLabel(&leadership.Member.Position)
			results = append(results, leadership)
		}
	}
	utils.Success(c, "", results)
}

func (h *RulerHandler) GetStateRulerLeadership(c *gin.Context) {
	stateId := c.Param("stateId")
	if stateId == "" {
		log.Errorf("StateID is required to get state ruler leadership")
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	results := make([]dto.LeadershipRes, 0)
	val, err := h.Service.GetStateLeadership(stateId)
	if err != nil {
		log.Errorf("Failed to fetch state leadership %+v", err)
		utils.Error(c, 500, "Failed to fetch state leadership", nil)
		return
	}
	for _, leader := range val {
		leadership := dto.LeadershipRes{}
		if err := utils.MapToStruct(leader, &leadership); err != nil {
			log.Errorf("Failed to parse state leadership: %+v ", err)
		} else {
			utils.TransformPositionLabel(&leadership.Position.Position)
			results = append(results, leadership)
		}
	}
	utils.Success(c, "", results)
}
