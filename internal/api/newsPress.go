package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func newsPressRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewNewsPressHandler(services)
	//make route with url param
	npStatement := rg.Group("/news-press")
	{
		// npStatement.GET("/", handler.ListState)
		npStatement.POST("/", handler.CreateNewsPress)
		npStatement.POST("/with-media", auth.AuthRequired(true), handler.CreateNewsPressWithMedia)
		npStatement.GET("/party/:partyID", handler.GetPostsByPartyID)
		npStatement.GET("/:id/detail", auth.AuthRequired(false), handler.GetDetailByID)
		npStatement.GET("/by-user", auth.AuthRequired(true), handler.GetByUserID)
	}
	return npStatement
}
