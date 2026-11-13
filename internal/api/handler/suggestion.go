package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/utils"
)

type SuggestionHandler struct {
	Service service.SuggestionService
}

func NewSuggestionHandler(services *service.Services) *SuggestionHandler {
	return &SuggestionHandler{
		Service: services.SuggestionService,
	}
}

func (h *SuggestionHandler) Create(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 500, "Failed to process request", err)
		return
	}

	var req dto.SuggestionReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req.CreatedBy = fmt.Sprintf("%d", userID)

	payload := dto.SuggestionCreate{
		SuggestionReq:   req,
		Tags:            parseSuggestionTags(c.PostForm("tags")),
		StatusRevisions: parseSuggestionStatusRevisions(c.PostForm("status_revisions")),
	}

	for i := 0; ; i++ {
		fileKey := fmt.Sprintf("documents[%d].file", i)
		files, ok := form.File[fileKey]
		if !ok {
			break
		}
		payload.Documents = append(payload.Documents, dto.SuggestionDocumentReq{
			File:     files[0],
			FileName: c.PostForm(fmt.Sprintf("documents[%d].filename", i)),
			Remark:   c.PostForm(fmt.Sprintf("documents[%d].remark", i)),
		})
	}

	resp, err := h.Service.Create(c, payload)
	if err != nil {
		utils.Error(c, 500, "Failed to create suggestion", err)
		return
	}
	utils.Success(c, "Suggestion submitted successfully", resp)
}

func (h *SuggestionHandler) ListByUser(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized access", nil)
		return
	}

	resp, err := h.Service.ListByUser(fmt.Sprintf("%d", userID))
	if err != nil {
		utils.Error(c, 500, "Failed to fetch suggestions", err)
		return
	}
	utils.Success(c, "Suggestions fetched successfully", resp)
}

func (h *SuggestionHandler) GetBySuggestionID(c *gin.Context) {
	suggestionID := c.Param("suggestionID")
	if suggestionID == "" {
		utils.Error(c, 400, "Suggestion ID is required", nil)
		return
	}

	resp, err := h.Service.GetBySuggestionID(suggestionID)
	if err != nil {
		utils.Error(c, 404, "Suggestion not found", err)
		return
	}
	utils.Success(c, "Suggestion fetched successfully", resp)
}

func parseSuggestionTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return nil
	}
	return tags
}

func parseSuggestionStatusRevisions(raw string) []dto.SuggestionStatusRevisionReq {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var revisions []dto.SuggestionStatusRevisionReq
	if err := json.Unmarshal([]byte(raw), &revisions); err != nil {
		return nil
	}
	return revisions
}
