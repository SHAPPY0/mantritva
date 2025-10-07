package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"gorm.io/gorm"
)

func vidhanSabhaRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewVidhansabhaRepository(db)
	service := service.NewVidhanSabhaService(repository)
	handler := handler.NewVidhanSabhaHandler(service)

	state := rg.Group("/vidhansabha")
	{
		state.GET("/", handler.ListVidhanSabha)
		state.POST("/", handler.CreateVidhanSabha)
	}
	return rg
}
