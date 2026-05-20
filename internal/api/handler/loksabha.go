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

type LoksabhaHandler struct {
	Service service.LoksabhaService
}

func NewLoksabhaHandler(serivce service.LoksabhaService) *LoksabhaHandler {
	return &LoksabhaHandler{
		Service: serivce,
	}
}

func (h *LoksabhaHandler) CreateLoksabha(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	var req dto.LoksabhaReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(authUser.UserName, req)
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

func (h *LoksabhaHandler) GetLoksabhaByStateId(c *gin.Context) {
	stateId := c.Param("stateId")
	if stateId == "" {
		utils.Error(c, 400, "State ID is required", nil)
		return
	}
	stateIdInt, err := utils.StrToInt(stateId)
	if err != nil {
		utils.Error(c, 400, "Invalid State ID", err)
		return
	}
	val, err := h.Service.GetByStateID(stateIdInt)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch loksabha", nil)
		return
	}
	utils.Success(c, "Loksabha fetched successfully", val)
}

func (h *LoksabhaHandler) CreateBulkLoksabha(c *gin.Context) {
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

	success, errored, err := h.Service.CreateBulkLoksabha(c, authUser.UserName, f)
	if err != nil {
		utils.Error(c, 500, "Failed to upload bulk loksabha", err)
		return
	}
	utils.Success(c, "Loksabha created successfully", map[string]any{
		"success": success,
		"errored": errored,
	})

}
