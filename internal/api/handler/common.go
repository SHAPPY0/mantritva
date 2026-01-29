package handler

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
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
