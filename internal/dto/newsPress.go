package dto

import (
	"mime/multipart"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsPressReq struct {
	Event         string                  `json:"event,omitempty" form:"event,omitempty"`
	MemberID      *uint                   `json:"member_id" form:"member_id,omitempty"`
	PartyID       *uint                   `json:"party_id" form:"party_id,omitempty"`
	EntityID      *uint                   `json:"entity_id" form:"entity_id,omitempty"`
	Subject       *utils.Subject          `json:"subject,omitempty" form:"subject,omitempty" binding:"omitempty"`
	Type          *utils.StatementType    `json:"type,omitempty" form:"type,omitempty" binding:"omitempty"`
	Tone          *utils.StatementTone    `json:"tone,omitempty" form:"tone,omitempty" binding:"omitempty"`
	Statement     string                  `json:"statement" binding:"required" form:"statement"`
	StatementDate utils.DateOnly          `json:"statement_date" binding:"required" form:"statement_date"`
	Status        utils.PostStatus        `json:"status" form:"status"`
	Source        *string                 `json:"source" binding:"required" form:"source"`
	Media         []*multipart.FileHeader `json:"media" form:"media"`
}

type NewsPressPostRes struct {
	model.NewsPressPost
	Member         Member           `json:"member,omitempty"`
	Entity         Entity           `json:"entity,omitempty"`
	Party          Party            `json:"party,omitempty"`
	Reaction       utils.Vote       `json:"reaction"`
	ReactionLocked bool             `json:"reaction_locked"`
	Questionnaire  *[]Questionnaire `json:"questionnaire,omitempty"`
	PostAdmin      *PostAdmin       `json:"post_admin,omitempty"`
}

type NewsPressStatus struct {
	ID     string           `json:"id"`
	Status utils.PostStatus `json:"status"`
}

type NewsPressShort struct {
	ID             primitive.ObjectID     `json:"id"`
	Event          *string                `json:"event"`
	Statement      string                 `json:"statement"`
	StatementDate  time.Time              `json:"statement_date"`
	Subject        string                 `json:"subject"`
	Type           string                 `json:"type"`
	Tone           string                 `json:"tone"`
	UpvoteCount    uint                   `json:"upvote_count"`
	DownvoteCount  uint                   `json:"downvote_count"`
	Member         Member                 `json:"member,omitempty"`
	IsEntity       bool                   `json:"is_entity"`
	Entity         Entity                 `json:"entity,omitempty"`
	Status         utils.PostStatus       `json:"status"`
	Party          Party                  `json:"party,omitempty"`
	Reaction       utils.Vote             `json:"reaction"`
	ReactionLocked bool                   `json:"reaction_locked"`
	Source         *string                `json:"source"`
	Media          []model.NewsPressMedia `json:"media,omitempty"`
	Questionnaire  *[]Questionnaire       `json:"questionnaire,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	PostAdmin      *PostAdmin             `json:"post_admin,omitempty"`
}
