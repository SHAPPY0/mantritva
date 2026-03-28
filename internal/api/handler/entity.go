package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type EntityHandler struct {
	Service service.EntityService
}

func NewEntityHandler(serivce service.EntityService) *EntityHandler {
	return &EntityHandler{
		Service: serivce,
	}
}

func (h *EntityHandler) CreateEntity(c *gin.Context) {
	var req dto.EntityReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create entity", err)
		return
	}
	utils.Success(c, "Entity created successfully", reps)
}

func (h *EntityHandler) ListEntity(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch entity", nil)
		return
	}
	utils.Success(c, "Entity fetched successfully", val)
}

func (h *EntityHandler) GetEntityByName(c *gin.Context) {
	name := c.Param("name")
	entityName := utils.NormalizeEntityName(name)
	val, err := h.Service.GetEntityByName(entityName)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch entity", nil)
		return
	}
	utils.Success(c, "", val)
}

func (h *EntityHandler) GetPostsByEntityID(c *gin.Context) {
	pagination := dto.NewPagination(c)
	entityID := c.Param("entityID")
	entityIdInt, err := utils.StrToInt(entityID)
	if err != nil {
		utils.Error(c, 400, "Invalid entity ID", nil)
		return
	}
	posts, err := h.Service.GetPostsByEntityID(pagination, entityIdInt)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch posts", err)
		return
	}
	utils.Success(c, "Posts fetched successfully", posts)
}
