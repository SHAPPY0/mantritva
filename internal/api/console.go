package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
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
		console.GET("/rulers/current", handler.GetCurrentRulersList)
		console.GET("/ruler/:rulerID/leadership", handler.GetLeadershipByRulerID)
		console.GET("/news-press", handler.GetNewsPressList)
		console.GET("/news-press/:id", handler.GetNewsPressPostByID)
		console.PUT("/news-press/:id", utils.AuthRequired(true), handler.UpdateNewsPressPost)
		console.GET("/grievances", handler.GetGrievancesList)
		console.GET("/grievance/:grievanceID", handler.GetGrievanceByID)
		console.PUT("/grievance/:grievanceID/updateStatus", handler.GrievanceUpdateStatus)
		console.POST("/leadership/exit/:id", handler.ExitLeadership)
	}
	return console
}
