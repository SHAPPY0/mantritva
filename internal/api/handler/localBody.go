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

type LocalBodyHandler struct {
	services *service.LocalBodyService
}

func NewLocalBodyHandler(services *service.Services) *LocalBodyHandler {
	return &LocalBodyHandler{
		services: services.LocalBodyService,
	}
}

func (h *LocalBodyHandler) CreateLocalBody(c *gin.Context) {
	username, err := auth.GetKeyInCtx[string](c, "user_name")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	if username == "" {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	var req dto.LocalBodyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Error while binding JSON ulb", err)
		return
	}
	if localBody, err := h.services.CreateLocalBody(username, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error while creating Local Body", err)
		return
	} else {
		utils.Success(c, "Local Body Created Successfully", localBody)
	}
}

func (h *LocalBodyHandler) BulkCreateLocalBodies(c *gin.Context) {
	username, err := auth.GetKeyInCtx[string](c, "user_name")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	if username == "" {
		utils.Error(c, 401, "Unauthorized Access", nil)
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

	success, errored, err := h.services.BulkCreateLocalBodies(username, c, f)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error while creating Local Bodies", err)
		return
	} else {
		utils.Success(c, "LocalBodies Created Successfully", map[string]any{
			"success": success,
			"errored": errored,
		})
	}
}

func (h *LocalBodyHandler) ListLocalBodiesByDistrict(c *gin.Context) {
	districtId := c.Param("districtId")
	if localBodies, err := h.services.ListLocalBodiesByDistrict(districtId); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error while fetching Local Bodies by district", err)
		return
	} else {
		result := make([]dto.LocalBodyRes, len(localBodies))
		for i, localBody := range localBodies {
			result[i] = dto.LocalBodyRes{
				LocalBody: localBody,
			}
			// result[i].LocalBodyType = utils.LOCAL_BODY_TYPE_URBAN
			result[i].DisplayType = localBody.BodyType.Label()
		}
		utils.Success(c, "LocalBodies fetched successfully", result)
	}
}
