package model

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
)

type Party struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Abbrv       string    `json:"abbrv"`
	Symbol      string    `json:"symbol"`
	About       string    `json:"about"`
	CoreBeliefs string    `json:"core_beliefs"`
	Leaning     string    `json:"leaning"`
	EcStatus    string    `json:"ec_status"`
	Founders    string    `json:"founders"`
	DOF         time.Time `json:"dof" gorm:"type:date"`
	AuditFields
}

type PartyMembers struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	PartyID       uint            `json:"party_id" gorm:"required"`
	MemberID      uint            `json:"member_id" gorm:"required"`
	Position      string          `json:"position" gorm:"required"`
	TermStartDate utils.DateOnly  `json:"term_start_date" gorm:"type:date"`
	TermEndDate   *utils.DateOnly `json:"term_end_date,omitempty" gorm:"type:date"`
	AuditFields
}
