package handler

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/utils"
)

type FeedHandler struct {
	Service                  service.FeedService
	ReportService            service.ReportService
	FeedReactionComputeQueue chan dto.ReactionReq
}

func NewFeedHandler(service *service.Services) *FeedHandler {
	h := &FeedHandler{
		Service:                  service.FeedService,
		ReportService:            service.ReportService,
		FeedReactionComputeQueue: make(chan dto.ReactionReq, 100),
	}
	return h
}

func (h *FeedHandler) GetLatestFeed(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	pagination := dto.NewPagination(c)
	feed, err := h.Service.GetLatestFeed(pagination, userID)
	if err != nil {
		utils.Error(c, 500, "Failed to load feed", err)
		return
	}
	utils.Success(c, "", feed)
}

func (h *FeedHandler) SaveReaction(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	var req dto.ReactionReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req.UserID = userID
	req.UserIP = clientIP(c)
	resp, err := h.Service.CreateReaction(req)
	if err != nil {
		utils.Error(c, 500, "Failed to react post", err)
		return
	}
	if resp {
		utils.Success(c, "", resp)
	} else {
		utils.Error(c, 500, "Failed to vote post", nil)
	}

}

// clientIP returns a normalized client IP (prefers X-Forwarded-For, converts IPv4-mapped IPv6).
func clientIP(c *gin.Context) string {
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	ip := c.ClientIP()
	p := net.ParseIP(ip)
	if p == nil {
		return ip
	}
	// IPv4-mapped IPv6 (::ffff:1.2.3.4) -> return IPv4
	if p4 := p.To4(); p4 != nil {
		return p4.String()
	}
	// map IPv6 loopback to IPv4 loopback
	if p.IsLoopback() {
		return "127.0.0.1"
	}
	return p.String()
}
