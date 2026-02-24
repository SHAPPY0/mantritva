package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func consoleRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	handler := handler.NewConsoleHandler(services)
	//make route with url param
	console := rg.Group("/console")
	{
		console.GET("/users", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetUsersList)
		console.PUT("/user/updateRole", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.UpdateRole)
		console.GET("/states", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetStatesList)
		console.GET("/parties", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetPartiesList)
		console.GET("/alliances", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetAlliancesList)
		console.GET("/members", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetMembersList)
		console.GET("/member/:id", auth.AuthRequired(false), auth.ValidateVersion(services.UserService), handler.GetMemberByID)
		console.PUT("/member/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.UpdateMemberByID)
		console.GET("/rulers", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetRulersList)
		console.GET("/rulers/current", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetCurrentRulersList)
		console.GET("/ruler/:rulerID/leadership", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetLeadershipByRulerID)
		console.GET("/news-press", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetNewsPressList)
		console.GET("/news-press/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetNewsPressPostByID)
		console.PUT("/news-press/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.UpdateNewsPressPost)
		console.GET("/grievances", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetGrievancesList)
		console.GET("/grievance/:grievanceID", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GetGrievanceByID)
		console.PUT("/grievance/:grievanceID/updateStatus", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.GrievanceUpdateStatus)
		console.POST("/leadership/exit/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), handler.ExitLeadership)
	}
	return console
}
