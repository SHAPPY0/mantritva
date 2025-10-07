package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
)

func leadershipRouter(rg *gin.RouterGroup, repos *repository.Repos) *gin.RouterGroup {
	service := service.NewLeadershipService(repos)
	handler := handler.NewLeadershipHandler(service)

	leadership := rg.Group("/leadership")
	{
		// leadership.GET("/", handler.List)
		leadership.POST("/", handler.Create)
	}
	return leadership
}
