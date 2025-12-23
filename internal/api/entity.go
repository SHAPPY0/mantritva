package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
)

func entityRouter(rg *gin.RouterGroup, repos *repository.Repos) *gin.RouterGroup {
	service := service.NewEntityService(repos)
	handler := handler.NewEntityHandler(service)

	entity := rg.Group("/entity")
	{
		entity.GET("/", handler.ListEntity)
		entity.POST("/", handler.CreateEntity)
		entity.GET("/:name", handler.GetEntityByName)
		entity.GET("/id/:entityID/news-press", handler.GetPostsByEntityID)
	}
	return rg
}
