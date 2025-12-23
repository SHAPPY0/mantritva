package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
)

func memberRouter(rg *gin.RouterGroup, repos *repository.Repos) *gin.RouterGroup {
	service := service.NewMemberService(repos)
	handler := handler.NewMemberHandler(service)

	member := rg.Group("/member")
	{
		member.GET("/", handler.ListMember)
		member.POST("/", handler.CreateMember)
		member.POST("/bulkUpload", handler.BulkMemberUpload)
		member.GET("/:name", handler.GetMemberByName)
		member.GET("/id/:memberID/career", handler.GetMemberCareer)
		member.GET("/id/:memberID/news-press", handler.GetPostsByMemberID)
		member.GET("/id/:memberID/grievances", handler.GetGrievancesByMemberID)
		member.GET("/unaffiliated", handler.ListUnaffiliatedMembers)
	}
	return rg
}
