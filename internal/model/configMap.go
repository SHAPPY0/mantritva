package model

type ConfigMap struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Category string `json:"category" gorm:"uniqueIndex:idx_category_key"`
	Key      string `json:"key" gorm:"uniqueIndex:idx_category_key"`
	Value    string `json:"value"`
	AuditFields
}
