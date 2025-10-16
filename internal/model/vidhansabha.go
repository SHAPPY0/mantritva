package model

type Vidhansabha struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	StateID    int    `json:"state_id"`
	DistrictID int    `json:"district_id"`
	ACNumber   int    `json:"ac_number"`
	AuditFields
}
