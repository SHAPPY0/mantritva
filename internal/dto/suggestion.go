package dto

import (
	"mime/multipart"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuggestionReq struct {
	Title           string                 `json:"title" form:"title" binding:"required"`
	Suggestion      string                 `json:"suggestion" form:"suggestion" binding:"required"`
	Topic           string                 `json:"topic" form:"topic" binding:"required"`
	Scope           string                 `json:"scope" form:"scope" binding:"required"`
	ScopeID         uint                   `json:"scope_id" form:"scope_id"`
	Status          utils.SuggestionStatus `json:"status" form:"status"`
	HideIdentity    bool                   `json:"hide_identity" form:"hide_identity"`
	RelatedArticles string                 `json:"related_articles" form:"related_articles"`
	CreatedBy       string                 `json:"created_by"`
}

type SuggestionDocumentReq struct {
	FileName string                `json:"filename" form:"filename"`
	Remark   string                `json:"remark" form:"remark"`
	File     *multipart.FileHeader `json:"file" form:"file"`
}

type SuggestionStatusRevisionReq struct {
	Date   string `json:"date"`
	Status string `json:"status"`
	Remark string `json:"remark"`
	Source string `json:"source"`
}

type SuggestionCreate struct {
	SuggestionReq
	Tags            []string
	StatusRevisions []SuggestionStatusRevisionReq
	Documents       []SuggestionDocumentReq
}

type SuggestionCreateRes struct {
	SuggestionID string `json:"suggestion_id"`
	Title        string `json:"title,omitempty"`
}

type SuggestionRes struct {
	model.Suggestion
	User UserClaim `json:"user,omitempty"`
}

type SuggestionStatusUpdate struct {
	ID     string                 `json:"id"`
	Date   utils.DateOnly         `json:"date"`
	Status utils.SuggestionStatus `json:"status"`
	Remark string                 `json:"remark"`
	Source string                 `json:"source"`
}

type SuggestionShort struct {
	ID             primitive.ObjectID     `json:"id"`
	SuggestionID   string                 `json:"suggestion_id,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Suggestion     string                 `json:"suggestion"`
	Topic          string                 `json:"topic,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Scope          string                 `json:"scope,omitempty"`
	ScopeID        *uint                  `json:"scope_id,omitempty"`
	Status         utils.SuggestionStatus `json:"status,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty"`
	Reaction       utils.Vote             `json:"reaction"`
	ReactionLocked bool                   `json:"reaction_locked"`
	UpvoteCount    uint                   `json:"upvote_count"`
	DownvoteCount  uint                   `json:"downvote_count"`
	PostAdmin      *PostAdmin             `json:"post_admin,omitempty"`
	Questionnaire  *[]Questionnaire       `json:"questionnaire,omitempty"`
}
