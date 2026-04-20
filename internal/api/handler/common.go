package handler

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type CommonHandler struct {
	Service service.CommonService
}

func NewCommonHandler(service *service.Services) *CommonHandler {
	return &CommonHandler{
		Service: service.CommonService,
	}
}

func (h *CommonHandler) ServeStatic(c *gin.Context) {
	objectKey := c.Param("objectkey")
	if objectKey == "" || objectKey == "/" {
		utils.Error(c, 400, "Object key is required", nil)
		return
	}
	filepath := strings.TrimPrefix(objectKey, "/")
	filepath = path.Clean(filepath)
	log.Infof("Serving static file: %s", filepath)
	url, err := h.Service.CreatePreSignedUrl(filepath)
	if err != nil {
		c.Data(http.StatusOK, "", []byte{})
		return
	}

	fileData, contentType, err := h.Service.DownloadFile(c, url)
	if err != nil {
		c.Data(http.StatusOK, "", []byte{})
		return
	}
	c.Data(http.StatusOK, contentType, *fileData)
}

func (h *CommonHandler) CreateEnquiry(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 400, "Invalid form data", err)
		return
	}
	var req dto.Enquiry
	if err := c.ShouldBind(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	files, _ := form.File["files"]

	ok, err := h.Service.CreateEnquiry(c, req, files)
	if err != nil {
		utils.Error(c, 500, "Internal Server Error", err)
		return
	}
	utils.Success(c, "", ok)
}

func (h *CommonHandler) WAWebhook(c *gin.Context) {
	var response map[any]any
	if err := c.ShouldBind(&response); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	log.Infof("WA webhook %v", response)
	utils.Success(c, "", nil)
}

func (h *CommonHandler) GetWAWebhook(c *gin.Context) {
	var response map[any]any
	if err := c.ShouldBind(&response); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	log.Infof("GET WA webhook %v", response)
	utils.Success(c, "", nil)
}
