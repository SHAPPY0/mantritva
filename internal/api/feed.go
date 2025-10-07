package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

func feedRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewFeedHandler(services)
	feed := rg.Group("/feed")
	{
		feed.GET("/", utils.AuthRequired(false), handler.GetLatestFeed)
		feed.POST("/reaction", utils.AuthRequired(true), handler.SaveReaction)
	}
	return feed
}
