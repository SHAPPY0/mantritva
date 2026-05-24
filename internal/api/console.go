package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func consoleRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	h := handler.NewConsoleHandler(services)
	botHandler := handler.NewBotHandler(services)

	console := rg.Group("/console")
	{
		console.GET("/users", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetUsersList)
		console.PUT("/user/updateRole", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.UpdateRole)
		console.POST("/user/addRole", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.AddRole)
		console.GET("/user/:userID/roles", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetUserRoles)
		console.GET("/role-scopes", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.ListRoleScopes)
		console.GET("/states", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetStatesList)
		console.GET("/parties", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetPartiesList)
		console.GET("/alliances", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetAlliancesList)
		console.GET("/members", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetMembersList)
		console.GET("/member/:id", auth.AuthRequired(false), auth.ValidateVersion(services.UserService), h.GetMemberByID)
		console.PUT("/member/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.UpdateMemberByID)
		console.GET("/rulers", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetRulersList)
		console.GET("/rulers/current", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetCurrentRulersList)
		console.GET("/ruler/:rulerID/leadership", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetLeadershipByRulerID)
		console.GET("/news-press", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetNewsPressList)
		console.GET("/news-press/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetNewsPressPostByID)
		console.PUT("/news-press/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.UpdateNewsPressPost)
		console.GET("/grievances", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetGrievancesList)
		console.GET("/grievance/:grievanceID", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetGrievanceByID)
		console.PUT("/grievance/:grievanceID/updateStatus", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GrievanceUpdateStatus)
		console.POST("/leadership/exit/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.ExitLeadership)
		console.GET("/enquiries", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetEnquiriesList)
		console.POST("/setup-bulk-leader", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.SetupBulkLeader)
	}
	bot := rg.Group("/bot")
	{
		bot.POST("/user/create", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), botHandler.CreateUser)
		bot.POST("/feed/react", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), botHandler.ReactFeeds)
	}
	return console
}
