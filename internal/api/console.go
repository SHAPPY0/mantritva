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
		console.GET("/users", utils.AuthRequired(true), handler.GetUsersList)
		console.GET("/states", utils.AuthRequired(true), handler.GetStatesList)
		console.GET("/parties", utils.AuthRequired(true), handler.GetPartiesList)
		console.GET("/alliances", utils.AuthRequired(true), handler.GetAlliancesList)
		console.GET("/members", utils.AuthRequired(true), handler.GetMembersList)
		console.GET("/member/:id", utils.AuthRequired(false), handler.GetMemberByID)
		console.PUT("/member/:id", utils.AuthRequired(true), handler.UpdateMemberByID)
		console.GET("/rulers", utils.AuthRequired(true), handler.GetRulersList)
		console.GET("/rulers/current", utils.AuthRequired(true), handler.GetCurrentRulersList)
		console.GET("/ruler/:rulerID/leadership", utils.AuthRequired(true), handler.GetLeadershipByRulerID)
		console.GET("/news-press", utils.AuthRequired(true), handler.GetNewsPressList)
		console.GET("/news-press/:id", utils.AuthRequired(true), handler.GetNewsPressPostByID)
		console.PUT("/news-press/:id", utils.AuthRequired(true), handler.UpdateNewsPressPost)
		console.GET("/grievances", utils.AuthRequired(true), handler.GetGrievancesList)
		console.GET("/grievance/:grievanceID", utils.AuthRequired(true), handler.GetGrievanceByID)
		console.PUT("/grievance/:grievanceID/updateStatus", utils.AuthRequired(true), handler.GrievanceUpdateStatus)
		console.POST("/leadership/exit/:id", utils.AuthRequired(true), handler.ExitLeadership)
	}
	return console
}
