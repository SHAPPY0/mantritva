package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func localBodyRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewLocalBodyHandler(services)

	localBody := rg.Group("/local-body")
	{
		localBody.POST("/create", auth.AuthRequired(true), handler.CreateLocalBody)
		localBody.POST("/bulk-create", auth.AuthRequired(true), handler.BulkCreateLocalBodies)
		localBody.GET("/by-district/:districtId", handler.ListLocalBodiesByDistrict)
	}
	return rg
}
