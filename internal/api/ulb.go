package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func ulbRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewULBHandler(services)

	ulb := rg.Group("/ulb")
	{
		ulb.POST("/create", auth.AuthRequired(true), handler.CreateULB)
		ulb.POST("/bulk-create", auth.AuthRequired(true), handler.BulkCreateULB)
		ulb.GET("/by-district/:districtId", handler.ListULBsByDistrict)
	}
	return rg
}
