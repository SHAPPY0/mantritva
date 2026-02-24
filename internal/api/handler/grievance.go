package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/utils"
)

type GrievanceHandler struct {
	Service service.GrievanceSerivce
}

func NewGrievanceHandler(services *service.Services) *GrievanceHandler {
	return &GrievanceHandler{
		Service: services.GrievanceService,
	}
}

func (h *GrievanceHandler) Create(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 500, "Failed to process request", err)
		return
	}
	var req dto.GrievanceReq
	var grievance dto.Grievance
	if err := c.ShouldBind(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req.CreatedBy = fmt.Sprintf("%d", userID)
	grievance.GrievanceReq = req
	for i := 0; ; i++ {
		var evidence dto.EvidenceReq
		fileKey := fmt.Sprintf("evidences[%d].file", i)
		files, ok := form.File[fileKey]
		if !ok {
			break
		}
		evidence.File = files[0]
		evidence.FileName = c.PostForm(fmt.Sprintf("evidences[%d].filename", i))
		evidence.Comment = c.PostForm(fmt.Sprintf("evidences[%d].comment", i))
		grievance.Evidences = append(grievance.Evidences, evidence)
	}

	resp, err := h.Service.Create(c, grievance)
	if err != nil {
		utils.Error(c, 500, "Failed to create grievance", err)
		return
	}
	utils.Success(c, "Grievance created successfully", resp)
}

func (h *GrievanceHandler) ListByUser(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized access", nil)
		return
	}

	resp, err := h.Service.ListByUser(fmt.Sprintf("%d", userID))
	if err != nil {
		utils.Error(c, 500, "Failed to fetch grievances", err)
		return
	}
	utils.Success(c, "Grievances fetched successfully", resp)
}

func (h *GrievanceHandler) GetByGrievanceID(c *gin.Context) {
	grievanceID := c.Param("grievanceID")
	if grievanceID == "" {
		utils.Error(c, 400, "Grievance ID is required", nil)
		return
	}

	resp, err := h.Service.GetByGrievanceID(grievanceID)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch grievance", err)
		return
	}
	utils.Success(c, "Grievance fetched successfully", resp)
}

func (h *GrievanceHandler) GetByStateID(c *gin.Context) {
	stateID := c.Param("stateID")
	if stateID == "" {
		utils.Error(c, 400, "StateID is required", nil)
		return
	}

	stateIDx, err := utils.StrToInt(stateID)
	if err != nil {
		utils.Error(c, 400, "Invalid stateID", nil)
		return
	}
	resp, err := h.Service.GetByStateID(stateIDx)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch grievances by stateID", err)
		return
	}
	utils.Success(c, "", resp)
}
