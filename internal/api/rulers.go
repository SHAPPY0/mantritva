package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
)

func rulersRouter(rg *gin.RouterGroup, repos *repository.Repos) *gin.RouterGroup {
	service := service.NewRulerService(repos)
	handler := handler.NewRulerHandler(service)

	ruler := rg.Group("/ruler")
	{
		ruler.GET("/", handler.ListRulers)
		ruler.POST("/", handler.CreateRuler)
		ruler.GET("/central", handler.GetCentralRuler)
		ruler.GET("/center/leadership", handler.GetCenterLeadership)
		ruler.GET("/state/:stateId", handler.GetStateRuler)
		ruler.GET("/state/:stateId/leadership", handler.GetStateRulerLeadership)
		ruler.GET("/rajyasabha/getByStateId/:stateId", handler.GetRajyasabhaByStateID)
	}
	return rg
}
