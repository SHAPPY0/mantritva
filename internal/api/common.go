package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
)

func commonRouter(rg *gin.RouterGroup, service *service.Services) *gin.RouterGroup {
	handler := handler.NewCommonHandler(service)
	common := rg.Group("/cmn")
	{
		common.GET("/serve-static/*objectkey", handler.ServeStatic)
	}
	return common
}
