package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/utils"
)

type ProtectedRoute struct {
	Path         string
	Method       string
	Roles        []utils.Roles
	AuthRequired bool
	Handler      gin.HandlerFunc
}

func consoleRouter(rg *gin.RouterGroup, services *service.Services) *gin.RouterGroup {
	h := handler.NewConsoleHandler(services)
	botHandler := handler.NewBotHandler(services)

	var protectedRoutes = []ProtectedRoute{
		{Path: "/users", Method: "GET", Handler: h.GetUsersList, Roles: []utils.Roles{utils.STATE_ADMIN}},
		{Path: "/states", Method: "GET", Handler: h.GetStatesByUser, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN}},
		{Path: "/districts", Method: "GET", Handler: h.GetDistrictsByUser, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN}},
		{Path: "/state/:stateId/districts", Method: "GET", Handler: h.GetDistrictsByStateByUser, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN}},
		{Path: "/district/:districtId/local-bodies", Method: "GET", Handler: h.GetLocalBodiesByDistrictByUser, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN, utils.LOCALBODY_ADMIN}},
		{Path: "/parties", Method: "GET", Handler: h.GetPartiesList},
		{Path: "/alliances", Method: "GET", Handler: h.GetAlliancesList},
		{Path: "/members", Method: "GET", Handler: h.GetMembersList},
		{Path: "/rulers", Method: "GET", Handler: h.GetRulersList, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN, utils.LOCALBODY_ADMIN}},
		{Path: "/rulers/current", Method: "GET", Handler: h.GetCurrentRulersList, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN, utils.LOCALBODY_ADMIN}},
		{Path: "/ruler/:rulerID/leadership", Method: "GET", Handler: h.GetLeadershipByRulerID, Roles: []utils.Roles{utils.STATE_ADMIN, utils.DISTRICT_ADMIN, utils.LOCALBODY_ADMIN}},
	}

	console := rg.Group("/console")
	{
		for _, protectedRoute := range protectedRoutes {
			authRequired := false
			if protectedRoute.AuthRequired {
				authRequired = true
			}
			switch protectedRoute.Method {
			case "GET":
				console.GET(protectedRoute.Path,
					auth.AuthRequired(authRequired),
					auth.ValidateVersion(services.UserService),
					auth.RequiredRoles(services.UserService, protectedRoute.Roles),
					protectedRoute.Handler)
			case "POST":
				console.POST(protectedRoute.Path,
					auth.AuthRequired(authRequired),
					auth.ValidateVersion(services.UserService),
					auth.RequiredRoles(services.UserService, protectedRoute.Roles),
					protectedRoute.Handler)
			case "PUT":
				console.PUT(protectedRoute.Path,
					auth.AuthRequired(authRequired),
					auth.ValidateVersion(services.UserService),
					auth.RequiredRoles(services.UserService, protectedRoute.Roles),
					protectedRoute.Handler)
			case "DELETE":
				console.DELETE(protectedRoute.Path,
					auth.AuthRequired(authRequired),
					auth.ValidateVersion(services.UserService),
					auth.RequiredRoles(services.UserService, protectedRoute.Roles),
					protectedRoute.Handler)
			}
		}
		// console.GET("/users", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetUsersList)
		console.GET("/dashboard", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetDashboard)
		console.PUT("/user/updateRole", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.UpdateRole)
		console.POST("/user/addRole", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.AddRole)
		console.GET("/user/:userID/roles", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetUserRoles)
		console.GET("/role-scopes", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.ListRoleScopes)
		// console.GET("/states", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetStatesList)
		// console.GET("/parties", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetPartiesList)
		// console.GET("/alliances", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetAlliancesList)
		// console.GET("/members", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetMembersList)
		console.GET("/member/:id", auth.AuthRequired(false), auth.ValidateVersion(services.UserService), h.GetMemberByID)
		console.PUT("/member/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.UpdateMemberByID)
		// console.GET("/rulers", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetRulersList)
		// console.GET("/rulers/current", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetCurrentRulersList)
		// console.GET("/ruler/:rulerID/leadership", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetLeadershipByRulerID)
		console.GET("/news-press", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetNewsPressList)
		console.GET("/news-press/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetNewsPressPostByID)
		console.PUT("/news-press/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.UpdateNewsPressPost)
		console.PUT("/news-press/:id/updateStatus", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.NewsPressUpdateStatus)
		console.GET("/grievances", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetGrievancesList)
		console.GET("/grievance/:grievanceID", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetGrievanceByID)
		console.PUT("/grievance/:grievanceID/updateStatus", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GrievanceUpdateStatus)
		console.GET("/development-welfares", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetDevelopmentWelfaresList)
		console.GET("/development-welfares/:id", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.GetDevelopmentWelfareByID)
		console.PUT("/development-welfares/:id/updateStatus", auth.AuthRequired(true), auth.ValidateVersion(services.UserService), h.DevelopmentWelfareUpdateStatus)
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
