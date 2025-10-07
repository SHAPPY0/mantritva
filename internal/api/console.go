package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
)

func consoleRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewConsoleHandler(services)
	//make route with url param
	console := rg.Group("/console")
	{
		console.GET("/users", handler.GetUsersList)
		console.GET("/states", handler.GetStatesList)
		console.GET("/parties", handler.GetPartiesList)
		console.GET("/alliances", handler.GetAlliancesList)
		console.GET("/members", handler.GetMembersList)
		console.GET("/rulers", handler.GetRulersList)
		console.GET("/ruler/:rulerID/leadership", handler.GetLeadershipByRulerID)
		console.GET("/news-press", handler.GetNewsPressList)
	}
	return console
}
