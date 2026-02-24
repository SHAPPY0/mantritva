package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func developmentWelfareRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewDevelopmentWelfareHandler(services)
	devWelfare := rg.Group("/development-welfare")
	{
		devWelfare.POST("/", auth.AuthRequired(true), handler.Create)
		devWelfare.GET("/", auth.AuthRequired(false), handler.List)
		devWelfare.GET("/:id", auth.AuthRequired(false), handler.GetByID)
		devWelfare.GET("/:id/detail", auth.AuthRequired(false), handler.GetDetailByID)
	}
	return devWelfare
}
