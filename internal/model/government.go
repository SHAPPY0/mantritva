package model

import (
	"time"
)

type Government struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PartyID    string    `json:"party_id"`
	AllianceID string    `json:"alliance_id"` // ID of the alliance if applicable
	Position   string    `json:"position"`    // e.g., "ruling", "opposition"
	Type       string    `json:"type"`        // e.g., "CENTRAL", "STATE"
	TermStart  time.Time `json:"term_start"`
	TermEnd    time.Time `json:"term_end"`
	IsAlliance bool      `json:"is_alliance"`
	IsCurrent  bool      `json:"is_current" gorm:"default:false"`
	AuditFields
}
