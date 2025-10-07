package model

import (
	"time"
)

type Alliance struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Abbrv      string `json:"abbrv"`
	FormedYear string `json:"formed_year"`
	IsActive   bool   `json:"is_active" gorm:"default:true"`
	Remark     string `json:"remark"` // Additional information about the alliance
	AuditFields
}

type AllianceParty struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AllianceID  int        `json:"alliance_id"`
	PartyID     int        `json:"party_id"`
	IsLeader    bool       `json:"is_leader"`
	JoiningDate *time.Time `json:"joining_date" gorm:"type:date"`
	LeavingDate *time.Time `json:"leaving_date,omitempty" gorm:"type:date"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	AuditFields
}
