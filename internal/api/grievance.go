package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func grievanceRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewGrievanceHandler(services)
	grievance := rg.Group("/grievance")
	{
		grievance.POST("/", auth.AuthRequired(true), handler.Create)
		grievance.POST("/:grievanceID/cosign", auth.AuthRequired(true), handler.Cosign)
		grievance.GET("/", auth.AuthRequired(true), handler.ListByUser)
		grievance.GET("/:grievanceID", auth.AuthRequired(false), handler.GetByGrievanceID)
		grievance.GET("/byStateID/:stateID", handler.GetByStateID)
	}
	return grievance
}
