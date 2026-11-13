package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func suggestionRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	h := handler.NewSuggestionHandler(services)
	suggestion := rg.Group("/public-suggestion")
	{
		suggestion.POST("/", auth.AuthRequired(true), h.Create)
		suggestion.GET("/", auth.AuthRequired(true), h.ListByUser)
		suggestion.GET("/:suggestionID", auth.AuthRequired(false), h.GetBySuggestionID)
	}
	return rg
}
