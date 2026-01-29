package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type DevelopmentWelfareHandler struct {
	Service     service.DevelopmentWelfareService
	FeedService service.FeedService
}

func NewDevelopmentWelfareHandler(services *service.Services) *DevelopmentWelfareHandler {
	return &DevelopmentWelfareHandler{
		Service:     services.DevelopmentWelfareService,
		FeedService: services.FeedService,
	}
}

func (h *DevelopmentWelfareHandler) Create(c *gin.Context) {
	userName, err := utils.GetKeyInCtx[string](c, "user_name")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 500, "Failed to process request", err)
		return
	}

	var req dto.DevelopmentWelfareReq
	// var grievance dto.Grievance
	if err := c.ShouldBind(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req.CreatedBy = userName

	for i := 0; ; i++ {
		var document dto.Document
		files, ok := form.File[fmt.Sprintf("documents[%d].file", i)]
		if !ok {
			break
		}
		document.File = files[0]
		document.FileName = c.PostForm(fmt.Sprintf("documents[%d].filename", i))
		document.Remark = c.PostForm(fmt.Sprintf("documents[%d].remark", i))
		log.Infof("Document %+v", document)
		req.Documents = append(req.Documents, document)
	}

	resp, err := h.Service.Create(c, req)
	if err != nil {
		utils.Error(c, 500, "Failed to create development welfare", err)
		return
	}

	feed := dto.FeedReq{
		Type:   utils.DEVELOPMENT,
		PostID: resp.ID,
		Status: utils.Published,
	}
	ok, err := h.FeedService.Create(feed)
	if err != nil {
		utils.Error(c, 500, "Failed to create development welfare feed record", err)
		return
	}
	if ok {
		utils.Success(c, "Development Welfare created successfully", resp)
	} else {
		utils.Error(c, 500, "Development Welfare created but error while feed creation", err)
	}

}

func (h *DevelopmentWelfareHandler) List(c *gin.Context) {
	resp, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to get development welfare list", err)
		return
	}
	utils.Success(c, "", resp)
}
