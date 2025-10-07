package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type NewsPressReq struct {
	MemberID      uint                `json:"member_id" binding:"required"`
	PartyID       uint                `json:"party_id" binding:"required"`
	Subject       utils.Subject       `json:"subject" binding:"required"`
	Type          utils.StatementType `json:"type" binding:"required"`
	Tone          utils.StatementTone `json:"tone" binding:"required"`
	Statement     string              `json:"statement" binding:"required"`
	StatementDate utils.DateOnly      `json:"statement_date" binding:"required"`
	Status        utils.PostStatus    `json:"status" binding:"required"`
}

type NewsPressPostRes struct {
	model.NewsPressPost
	Member   Member                 `json:"member"`
	Party    Party                  `json:"party"`
	Reaction utils.PostReactionType `json:"reaction"`
}
