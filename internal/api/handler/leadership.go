package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LeadershipHandler struct {
	Service service.LeadershipService
}

func NewLeadershipHandler(service service.LeadershipService) *LeadershipHandler {
	return &LeadershipHandler{
		Service: service,
	}
}

func (h *LeadershipHandler) Create(c *gin.Context) {
	var req dto.LeadershipReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	resp, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create leadership", err)
		return
	}
	utils.Success(c, "Leadership created successfully", resp)
}

func (h *LeadershipHandler) CreateBulkLeaders(c *gin.Context) {
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
	success, errored, err := h.Service.CreateBulk(f)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Error while processing bulk leadership create", err)
		return
	}
	utils.Success(c, "Bulk leaders created successfully", map[string]any{
		"success": success,
		"errored": errored,
	})
}
