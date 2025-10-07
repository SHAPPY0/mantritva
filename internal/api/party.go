package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
)

func partyRouter(rg *gin.RouterGroup, repos *repository.Repos) *gin.RouterGroup {
	service := service.NewPartyService(repos)
	handler := handler.NewPartyHandler(service)

	party := rg.Group("/party")
	{
		party.GET("/", handler.ListParty)
		party.POST("/", handler.CreateParty)
		party.GET("/:abbrv", handler.GetPartyByAbbrv)
		party.PUT("/:partyId/uploadSymbol", handler.UploadPartySymbol)
		party.POST("/addMember", handler.AddMember)
		party.GET("/members/:partyId", handler.GetPartyMembers)
		party.GET("/suggestion", handler.GetRandomParties)
	}
	return rg
}
