package model

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suggestion struct {
	ID                 primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	SuggestionID       string                 `json:"suggestion_id" bson:"suggestion_id"`
	Title              string                 `json:"title" bson:"title"`
	Suggestion         string                 `json:"suggestion" bson:"suggestion"`
	Documents          []Document             `json:"documents" bson:"documents"`
	Topic              string                 `json:"topic" bson:"topic"`
	Tags               []string               `json:"tags" bson:"tags"`
	Scope              string                 `json:"scope" bson:"scope"`
	ScopeID            *uint                  `json:"scope_id,omitempty" bson:"scope_id,omitempty"`
	RelatedArticles    []string               `json:"related_articles,omitempty" bson:"related_articles,omitempty"`
	HideIdentity       bool                   `json:"hide_identity" bson:"hide_identity"`
	Status             utils.SuggestionStatus `json:"status" bson:"status"`
	StatusRevisions    []StatusRevisions      `json:"status_revisions" bson:"status_revisions"`
	GovernmentResponse *GovernmentResponse    `json:"government_response,omitempty" bson:"government_response,omitempty"`
	UpvoteCount        uint                   `json:"upvote_count" bson:"upvote_count"`
	DownvoteCount      uint                   `json:"downvote_count" bson:"downvote_count"`
	ViewCount          uint                   `json:"view_count" bson:"view_count"`
	AssignedEntityID   *uint                  `json:"assigned_entity_id,omitempty" bson:"assigned_entity_id,omitempty"`
	AdminNotes         string                 `json:"admin_notes,omitempty" bson:"admin_notes,omitempty"`
	IsFeatured         bool                   `json:"is_featured" bson:"is_featured"`
	MergedIntoID       *primitive.ObjectID    `json:"merged_into_id,omitempty" bson:"merged_into_id,omitempty"`
	AuditFieldsBSON    `bson:",inline"`
}

type GovernmentResponse struct {
	Text        string    `json:"text" bson:"text"`
	RespondedBy string    `json:"responded_by" bson:"responded_by"`
	RespondedAt time.Time `json:"responded_at" bson:"responded_at"`
	Status      string    `json:"status" bson:"status"`
}
