package dto

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
)

type PartyReq struct {
	Name        string    `json:"name" binding:"required"`
	Abbrv       string    `json:"abbrv" binding:"required"`
	Symbol      string    `json:"symbol"`
	About       string    `json:"about"`
	CoreBeliefs string    `json:"core_beliefs"`
	Leaning     string    `json:"leaning"`
	EcStatus    string    `json:"ec_status"`
	Founders    string    `json:"founders"`
	DOF         time.Time `json:"dof"`
}

type AddMemberReq struct {
	PartyId   uint           `json:"party_id" binding:"required"`
	MemberId  uint           `json:"member_id" binding:"required"`
	Position  string         `json:"position"`
	TermStart utils.DateOnly `json:"term_start"`
	TermEnd   utils.DateOnly `json:"term_end,omitempty"`
}

type PartyMemberRes struct {
	ID       uint   `json:"id"`
	MemberID uint   `json:"member_id"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	Position string `json:"position"`
}

type Party struct {
	Name   string `json:"name"`
	Abbrv  string `json:"abbrv"`
	Symbol string `json:"symbol"`
}
