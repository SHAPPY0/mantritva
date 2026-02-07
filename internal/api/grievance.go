package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

func grievanceRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewGrievanceHandler(services)
	grievance := rg.Group("/grievance")
	{
		grievance.POST("/", utils.AuthRequired(true), handler.Create)
		grievance.GET("/", utils.AuthRequired(true), handler.ListByUser)
		grievance.GET("/:grievanceID", utils.AuthRequired(false), handler.GetByGrievanceID)
	}
	return grievance
}
