package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

func developmentWelfareRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewDevelopmentWelfareHandler(services)
	devWelfare := rg.Group("/development-welfare")
	{
		devWelfare.POST("/", utils.AuthRequired(true), handler.Create)
		devWelfare.GET("/", utils.AuthRequired(false), handler.List)
		devWelfare.GET("/:id", utils.AuthRequired(false), handler.GetByID)
	}
	return devWelfare
}
