package model

type Country struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Abbrv   string `json:"abbrv"`
	Icon    string `json:"icon"`
	Capital string `json:"capital"`
	AuditFields
}
