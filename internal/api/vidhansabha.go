package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"gorm.io/gorm"
)

func vidhanSabhaRouter(rg *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	repository := repository.NewVidhansabhaRepository(db)
	service := service.NewVidhanSabhaService(repository)
	handler := handler.NewVidhanSabhaHandler(service)

	vs := rg.Group("/vidhansabha")
	{
		vs.GET("/", handler.ListVidhanSabha)
		vs.POST("/", auth.AuthRequired(true), handler.CreateVidhanSabha)
		vs.POST("/bulk-create", auth.AuthRequired(true), handler.CreateBulkVidhanSabha)
		vs.GET("/getByStateId/:id", handler.GetByStateID)

	}
	return rg
}
