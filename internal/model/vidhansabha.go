package model

type Vidhansabha struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	LoksabhaID int    `json:"loksabha_id"`
	ACNumber   int    `json:"ac_number"`
	AuditFields
}
