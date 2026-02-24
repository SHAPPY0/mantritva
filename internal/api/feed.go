package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func feedRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewFeedHandler(services)
	feed := rg.Group("/feed")
	{
		feed.GET("/", auth.AuthRequired(false), handler.GetLatestFeed)
		feed.POST("/reaction", auth.AuthRequired(true), handler.SaveReaction)
	}
	return feed
}
