package model

import (
	"github.com/shappy0/mantritva/pkg/utils"
)

type Ruler struct {
	ID         uint                 `json:"id" gorm:"primaryKey"`
	Scope      utils.RulerScopeType `json:"scope"`
	StateID    int                  `json:"state_id"`
	IsCentral  bool                 `json:"is_central"`
	PartyID    int                  `json:"party_id"`
	AllianceID int                  `json:"alliance_id"`
	IsAlliance bool                 `json:"is_alliance" gorm:"default:false"`
	StartDate  utils.DateOnly       `json:"start_date" gorm:"type:date"`
	EndDate    utils.DateOnly       `json:"end_date" gorm:"type:date"`
	IsActive   bool                 `json:"is_active" gorm:"default:true"`
	AuditFields
}
