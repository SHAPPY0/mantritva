package dto

import (
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsPressReq struct {
	Event         string              `json:"event,omitempty"`
	MemberID      uint                `json:"member_id"`
	PartyID       uint                `json:"party_id"`
	EntityID      uint                `json:"entity_id"`
	Subject       utils.Subject       `json:"subject" binding:"required"`
	Type          utils.StatementType `json:"type" binding:"required"`
	Tone          utils.StatementTone `json:"tone" binding:"required"`
	Statement     string              `json:"statement" binding:"required"`
	StatementDate utils.DateOnly      `json:"statement_date" binding:"required"`
	Status        utils.PostStatus    `json:"status" binding:"required"`
	Source        *string             `json:"source"`
}

type NewsPressPostRes struct {
	model.NewsPressPost
	Member        Member           `json:"member,omitempty"`
	Entity        Entity           `json:"entity,omitempty"`
	Party         Party            `json:"party,omitempty"`
	Reaction      utils.Vote       `json:"reaction"`
	Questionnaire *[]Questionnaire `json:"questionnaire,omitempty"`
}

type NewsPressShort struct {
	ID             primitive.ObjectID `json:"id"`
	Event          *string            `json:"event"`
	Statement      string             `json:"statement"`
	StatementDate  time.Time          `json:"statement_date"`
	Subject        string             `json:"subject"`
	Type           string             `json:"type"`
	Tone           string             `json:"tone"`
	UpvoteCount    uint               `json:"upvote_count"`
	DownvoteCount  uint               `json:"downvote_count"`
	Member         Member             `json:"member,omitempty"`
	Entity         Entity             `json:"entity,omitempty"`
	Party          Party              `json:"party,omitempty"`
	Reaction       utils.Vote         `json:"reaction"`
	ReactionLocked bool               `json:"reaction_locked"`
	Source         *string            `json:"source"`
	Questionnaire  *[]Questionnaire   `json:"questionnaire,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
}
