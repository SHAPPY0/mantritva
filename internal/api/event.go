package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func eventRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewEventHandler(services)
	event := rg.Group("/event")
	{
		event.POST("/", auth.AuthRequired(true), handler.Create)
		event.GET("/", auth.AuthRequired(true), handler.GetAll)
	}
	return event
}
