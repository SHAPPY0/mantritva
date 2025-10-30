package model

type Entity struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Abbrv      string `json:"abbrv"`
	IsCentral  bool   `json:"is_central"`
	StateID    int    `json:"state_id"`
	DistrictID int    `json:"district_id"`
	AuditFields
}
