package model

type Loksabha struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	StateID    int    `json:"state_id"`
	DistrictID int    `json:"district_id"`
	PCNumber   int    `json:"pc_number"`
	AuditFields
}
