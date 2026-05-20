package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
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
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	var req dto.VidhanReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(authUser.UserName, req)
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

func (h *VidhanSabhaHandler) CreateBulkVidhanSabha(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to get uploaded file", err)
		return
	}

	f, err := fh.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to open uploaded file", err)
		return
	}
	defer f.Close()

	ext := filepath.Ext(fh.Filename)
	if ext != ".json" {
		utils.Error(c, http.StatusBadRequest, "Only .json file accepted", err)
		return
	}

	success, errored, err := h.Service.CreateBulkVidhansabha(c, authUser.UserName, f)
	if err != nil {
		utils.Error(c, 500, "Failed to upload bulk vidhansabha", err)
		return
	}
	utils.Success(c, "Vidhansabha created successfully", map[string]any{
		"success": success,
		"errored": errored,
	})
}
