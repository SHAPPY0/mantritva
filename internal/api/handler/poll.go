package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type PollHandler struct{ Service service.PollService }

func NewPollHandler(s *service.Services) *PollHandler { return &PollHandler{Service: s.PollService} }
func (h *PollHandler) Create(c *gin.Context) {
	uid, e := auth.GetKeyInCtx[int](c, "user_id")
	if e != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	var req dto.PollCreateReq
	if e = c.ShouldBindJSON(&req); e != nil {
		utils.Error(c, 400, "Invalid poll data", e)
		return
	}
	r, e := h.Service.Create(uid, req)
	if e != nil {
		utils.Error(c, 400, e.Error(), nil)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": r})
}
func (h *PollHandler) Get(c *gin.Context) {
	id, e := primitive.ObjectIDFromHex(c.Param("pollID"))
	if e != nil {
		utils.Error(c, 400, "Invalid poll ID", e)
		return
	}
	uid, _ := auth.GetKeyInCtx[int](c, "user_id")
	r, e := h.Service.Get(uid, id)
	if e != nil {
		utils.Error(c, 404, "Poll not found", e)
		return
	}
	utils.Success(c, "", r)
}
func (h *PollHandler) Vote(c *gin.Context) {
	uid, e := auth.GetKeyInCtx[int](c, "user_id")
	if e != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	id, e := primitive.ObjectIDFromHex(c.Param("pollID"))
	if e != nil {
		utils.Error(c, 400, "Invalid poll ID", e)
		return
	}
	var req dto.PollVoteReq
	if e = c.ShouldBindJSON(&req); e != nil {
		utils.Error(c, 400, "Invalid vote data", e)
		return
	}
	if e = h.Service.Vote(uid, id, req.OptionID); e != nil {
		utils.Error(c, 400, e.Error(), nil)
		return
	}
	utils.Success(c, "Vote recorded", nil)
}
