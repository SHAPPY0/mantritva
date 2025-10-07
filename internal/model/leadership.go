package model

import (
	"github.com/shappy0/mantritva/pkg/utils"
)

type Leadership struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	RulerID     uint            `json:"ruler_id"`
	MemberID    uint            `json:"member_id"`
	PartyID     uint            `json:"party_id"`
	Position    string          `json:"position"`
	Role        string          `json:"role,omitempty"`
	Departments string          `json:"departments,omitempty"`
	StateID     uint            `json:"state_id"`
	LoksabhaID  uint            `json:"loksabha_id,omitempty"`
	VidhansabhaID uint            `json:"vidhansabha_id,omitempty"`
	RajyasabhaID uint            `json:"rajyasabha_id,omitempty"`
	IsCentral   bool            `json:"is_central"`
	TermStart   utils.DateOnly  `json:"term_start" gorm:"type:date"`
	TermEnd     *utils.DateOnly `json:"term_end" gorm:"type:date"`
	IsActive    bool            `json:"is_active" gorm:"default:true"`
	AuditFields
}
