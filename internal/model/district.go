package model

type District struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	StateID int    `json:"state_id"`
	Abbrv   string `json:"abbrv"`
	AuditFields
}
