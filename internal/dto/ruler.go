package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type RulerReq struct {
	StateID    int            `json:"state_id"`
	IsCentral  bool           `json:"is_central"`
	PartyID    int            `json:"party_id" required:"true"`
	AllianceID int            `json:"alliance_id"`
	IsAlliance bool           `json:"is_alliance"`
	StartDate  utils.DateOnly `json:"start_date" required:"true"`
	EndDate    utils.DateOnly `json:"end_date,omitempty"`
	IsActive   bool           `json:"is_active"`
}

type RulerRes struct {
	model.Ruler
	Party      Party              `json:"party"`
	Leaders    []Leader           `json:"leaders"`
	IsAlliance bool               `json:"is_alliance"`
	Alliance   AllianceName       `json:"alliance"`
	State      *model.StateDetail `json:"state,omitempty"`
}
