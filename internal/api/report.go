package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
)

func reportRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewReportHandler(services)

	report := rg.Group("/report")
	{
		report.GET("/member/:memberID/generate", handler.GenerateReport)
		report.GET("/member/:memberID", handler.GetReport)
		report.GET("/:name/download", handler.Download)
	}
	return rg
}
