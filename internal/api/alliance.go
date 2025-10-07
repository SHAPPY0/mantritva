package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"gorm.io/gorm"
)

func allianceRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewAllianceRepository(db)
	service := service.NewAllianceService(repository)
	handler := handler.NewAllianceHandler(service)

	alliance := rg.Group("/alliance")
	{
		alliance.GET("/", handler.ListAlliance)
		alliance.POST("/", handler.CreateAlliance)
		alliance.POST("/party", handler.AddPartyToAlliance)
		alliance.GET("/:alliance", handler.GetAllianceByAbbrv)
	}
	return rg
}
