package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"gorm.io/gorm"
)

func stateRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewStateRepository(db)
	service := service.NewStateService(repository)
	handler := handler.NewStateHandler(service)
	//make route with url param
	state := rg.Group("/state")
	{
		state.GET("/", handler.ListState)
		state.POST("/", handler.CreateState)
		state.GET("/:name", handler.GetState)
	}
	return rg
}
