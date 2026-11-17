package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/api/handler"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
)

func pollRouter(rg *gin.RouterGroup, s *service.Services) *gin.RouterGroup {
	h := handler.NewPollHandler(s)
	g := rg.Group("/polls")
	g.POST("", auth.AuthRequired(true), h.Create)
	g.GET("/:pollID", auth.AuthRequired(false), h.Get)
	g.POST("/:pollID/vote", auth.AuthRequired(true), h.Vote)
	return g
}
