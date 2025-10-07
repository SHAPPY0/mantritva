package dto

import (
	"time"

	"github.com/shappy0/mantritva/internal/model"
)

type AllianceReq struct {
	Name       string `json:"name" binding:"required"`
	Abbrv      string `json:"abbrv" binding:"required"`
	IsActive   bool   `json:"is_active"`
	FormedYear string `json:"formed_year" binding:"required"`
	Remark     string `json:"remark"`
}

type AlliancePartyReq struct {
	AllianceID  int        `json:"alliance_id" binding:"required"`
	PartyID     int        `json:"party_id" binding:"required"`
	IsLeader    bool       `json:"is_leader"`
	JoiningDate *time.Time `json:"joining_date"`
	LeavingDate *time.Time `json:"leaving_date,omitempty"`
	IsActive    bool       `json:"is_active" binding:"required"`
}

type AllianceName struct {
	Name  string `json:"name"`
	Abbrv string `json:"abbrv"`
}

type AllianceParty struct {
	model.AllianceParty
	Party model.Party `json:"party"`
}

type AlliancePartyDetail struct {
	model.Alliance
	AllianceParty []AllianceParty `json:"alliance_parties"`
}
