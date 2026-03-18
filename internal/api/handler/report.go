package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/log"
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

	utils.Success(c, "", report)
}

func (h *ReportHandler) Download(c *gin.Context) {
	url := c.Query("url")
	log.Infof("URL %s", url)
	if url == "" {
		utils.Error(c, 404, "Please provide report url to download", nil)
		return
	}
	reportPdf, err := h.Service.DownloadReport(url, "report-card")
	if err != nil {
		utils.Error(c, 501, "Error downloading report", err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="report-%s.pdf"`, "memberID"))
	c.Data(http.StatusOK, "application/pdf", reportPdf)
}
