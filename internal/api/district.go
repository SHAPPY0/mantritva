package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"gorm.io/gorm"
)

func districtRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewDistrictRepository(db)
	service := service.NewDistrictService(repository)
	handler := handler.NewDistrictHandler(service)

	district := rg.Group("/district")
	{
		district.GET("/", handler.ListDistrict)
		district.POST("/", handler.CreateDistrict)
		district.POST("/bulk", handler.CreateBulkDistrict)
		district.GET("/getByStateId/:stateId", handler.GetDistrictByState)
	}
	return rg
}
