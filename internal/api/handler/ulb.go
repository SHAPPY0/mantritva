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

type UlbHandler struct {
	services *service.UlbService
}

func NewULBHandler(services *service.Services) *UlbHandler {
	return &UlbHandler{
		services: services.UlbService,
	}
}

func (h *UlbHandler) CreateULB(c *gin.Context) {
	username, err := auth.GetKeyInCtx[string](c, "user_name")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	if username == "" {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	var req dto.ULBReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Error while binding JSON ulb", err)
		return
	}
	if ulb, err := h.services.CreateULB(username, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error while creating ULB", err)
		return
	} else {
		utils.Success(c, "ULB Created Successfully", ulb)
	}
}

func (h *UlbHandler) BulkCreateULB(c *gin.Context) {
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

	success, errored, err := h.services.BulkCreateULB(username, c, f)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error while creating ULBs", err)
		return
	} else {
		utils.Success(c, "ULBs Created Successfully", map[string]any{
			"success": success,
			"errored": errored,
		})
	}
}

func (h *UlbHandler) ListULBsByDistrict(c *gin.Context) {
	districtId := c.Param("districtId")
	if ulbs, err := h.services.ListULBsByDistrict(districtId); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Error while fetching ULBs by district", err)
		return
	} else {
		result := make([]dto.LocalBodyRes, len(ulbs))
		for i, ulb := range ulbs {
			result[i] = dto.LocalBodyRes{
				ULB: ulb,
			}
			result[i].LocalBodyType = utils.LOCAL_BODY_TYPE_URBAN
			result[i].DisplayType = ulb.BodyType.Label()
		}
		utils.Success(c, "ULBs fetched successfully", result)
	}
}
