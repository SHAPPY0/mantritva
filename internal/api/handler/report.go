package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type ReportHandler struct {
	Service service.ReportService
}

func NewReportHandler(s *service.Services) *ReportHandler {
	return &ReportHandler{
		Service: s.ReportService,
	}
}

func (h *ReportHandler) GenerateReport(c *gin.Context) {
	memberID, err := utils.StrToInt(c.Param("memberID"))
	if err != nil {
		utils.Error(c, 404, "Invlid memberID", nil)
		return
	}
	report, err := h.Service.GenerateReport(memberID)
	if err != nil {
		utils.Error(c, 501, "Error getting report", err)
		return
	}
	utils.Success(c, "", report)
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	memberID, err := utils.StrToInt(c.Param("memberID"))
	if err != nil {
		utils.Error(c, 404, "Invlid memberID", nil)
		return
	}
	report, err := h.Service.GetReport(memberID)
	if err != nil {
		utils.Error(c, 501, "Error getting report", err)
		return
	}
	if report == nil {
		utils.Success(c, "No report found", nil)
		return
	}
	utils.Success(c, "", report)
}
