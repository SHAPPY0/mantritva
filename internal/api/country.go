package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"gorm.io/gorm"
)

func countryRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewCountryRepository(db)
	service := service.NewCountryService(repository)
	handler := handler.NewCountryHandler(service)

	user := rg.Group("/country")
	{
		user.GET("/", handler.ListCountry)
		user.POST("/", handler.CreateCountry)
	}
	return rg
}
