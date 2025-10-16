package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"gorm.io/gorm"
)

func loksabhaRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewLoksabhaRepository(db)
	service := service.NewLoksabhaService(repository)
	handler := handler.NewLoksabhaHandler(service)

	ls := rg.Group("/loksabha")
	{
		ls.GET("/", handler.ListLoksabha)
		ls.POST("/", handler.CreateLoksabha)
		ls.GET("/getByStateId/:stateId", handler.GetLoksabhaByStateId)
	}
	return rg
}
