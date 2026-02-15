package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

func newsPressRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewNewsPressHandler(services)
	//make route with url param
	npStatement := rg.Group("/news-press")
	{
		// npStatement.GET("/", handler.ListState)
		npStatement.POST("/", handler.CreateNewsPress)
		npStatement.GET("/party/:partyID", handler.GetPostsByPartyID)
		npStatement.GET("/:id/detail", utils.AuthRequired(false), handler.GetDetailByID)
	}
	return npStatement
}
