package service

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuggestionService interface {
	Create(*gin.Context, dto.SuggestionCreate) (*dto.SuggestionCreateRes, error)
	ListByUser(string) ([]model.Suggestion, error)
	GetBySuggestionID(string) (*model.Suggestion, error)
}

type suggestionService struct {
	Repository repository.SuggestionRepository
}

func NewSuggestionService(repos *repository.Repos) SuggestionService {
	return &suggestionService{
		Repository: repos.Suggestion,
	}
}

func (s *suggestionService) Create(ctx *gin.Context, req dto.SuggestionCreate) (*dto.SuggestionCreateRes, error) {
	if req.Scope != "NATIONAL" && req.ScopeID == 0 {
		return nil, utils.NewError("scope_id is required for non-national suggestions")
	}

	status := req.Status
	if status == "" {
		status = utils.SUGGESTION_SUBMITTED
	}

	suggestionID := generateSuggestionID(5)
	suggestion := &model.Suggestion{
		ID:              primitive.NewObjectID(),
		SuggestionID:    suggestionID,
		Title:           strings.TrimSpace(req.Title),
		Suggestion:      req.Suggestion,
		Documents:       make([]model.Document, 0),
		Topic:           req.Topic,
		Tags:            req.Tags,
		Scope:           req.Scope,
		RelatedArticles: parseRelatedArticles(req.RelatedArticles),
		HideIdentity:    req.HideIdentity,
		Status:          status,
		StatusRevisions: buildStatusRevisions(req.StatusRevisions, string(status)),
		AuditFieldsBSON: model.InsertBSONAuditFields(req.CreatedBy),
	}

	if req.ScopeID > 0 {
		scopeID := req.ScopeID
		suggestion.ScopeID = &scopeID
	}

	destinationPath := filepath.Join("suggestions", suggestion.SuggestionID)
	fileStorage, err := storage.New(destinationPath)
	if err != nil {
		log.Infof("Error initiating storage connection %v", err)
		return nil, err
	}

	for _, documentReq := range req.Documents {
		if documentReq.File == nil {
			continue
		}
		fileMeta, err := fileStorage.Save(ctx, documentReq.File)
		if err != nil {
			log.Warnf("Error while storing suggestion document %v", err)
			continue
		}
		suggestion.Documents = append(suggestion.Documents, model.Document{
			Url:      fileMeta.DownloadURL,
			FileName: fileMeta.FullPath,
			FileType: fileMeta.Extenstion,
			Remark:   documentReq.Remark,
		})
	}

	if err := s.Repository.Create(suggestion); err != nil {
		fileStorage.DeleteAll(suggestion.Documents)
		return nil, err
	}

	return &dto.SuggestionCreateRes{
		SuggestionID: suggestion.SuggestionID,
		Title:        suggestion.Title,
	}, nil
}

func (s *suggestionService) ListByUser(userID string) ([]model.Suggestion, error) {
	return s.Repository.ListByUser(userID)
}

func (s *suggestionService) GetBySuggestionID(suggestionID string) (*model.Suggestion, error) {
	suggestion, err := s.Repository.GetBySuggestionID(suggestionID)
	if err != nil {
		return nil, err
	}
	if err := s.Repository.IncrementViewCount(suggestion.ID); err != nil {
		log.Warnf("Failed to increment view count for suggestion %s: %v", suggestionID, err)
	}
	return suggestion, nil
}

func generateSuggestionID(n int) string {
	if n <= 0 {
		return ""
	}
	var sb strings.Builder
	first, _ := rand.Int(rand.Reader, big.NewInt(9))
	sb.WriteByte(byte('1' + first.Int64()))
	for i := 1; i < (n - 1); i++ {
		d, _ := rand.Int(rand.Reader, big.NewInt(10))
		sb.WriteByte(byte('0' + d.Int64()))
	}
	year := time.Now().Format("06")
	return "ps" + year + sb.String()
}

func parseRelatedArticles(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var articles []string
	if err := json.Unmarshal([]byte(raw), &articles); err == nil {
		return articles
	}
	return []string{raw}
}

func buildStatusRevisions(
	revisions []dto.SuggestionStatusRevisionReq,
	defaultStatus string,
) []model.StatusRevisions {
	if len(revisions) == 0 {
		return []model.StatusRevisions{
			{
				Date:   time.Now(),
				Status: defaultStatus,
				Remark: "Suggestion submitted by citizen",
				Source: "citizen",
			},
		}
	}

	result := make([]model.StatusRevisions, 0, len(revisions))
	for _, revision := range revisions {
		revisionDate := time.Now()
		if revision.Date != "" {
			if parsed, err := time.Parse(time.RFC3339, revision.Date); err == nil {
				revisionDate = parsed
			}
		}
		result = append(result, model.StatusRevisions{
			Date:   revisionDate,
			Status: revision.Status,
			Remark: revision.Remark,
			Source: revision.Source,
		})
	}
	return result
}
