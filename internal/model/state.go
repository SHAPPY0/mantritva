package model

type State struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Slug      string `json:"slug" gorm:"uniqueIndex"`
	CountryID int    `json:"country_id"`
	Abbrv     string `json:"abbrv"`
	Capital   string `json:"capital"`
	Structure string `json:"structure"`
	AuditFields
}

type StateDetail struct {
	State
	// VidhanSabhas []Vidhansabha `json:"vidhanSabhas"`
	// RayjaSabhas []Rayjasabha `json:"rayjaSabhas"`
	// LokSabhas []Loksabha `json:"lokSabhas"`
}
