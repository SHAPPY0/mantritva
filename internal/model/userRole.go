package model

type UserRole struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      string `json:"user_id" gorm:"uniqueIndex:idx_user_role"`
	RoleID      uint   `json:"role_id" gorm:"uniqueIndex:idx_user_role"`
	StateID     uint   `json:"state_id"`
	DistrictID  uint   `json:"district_id"`
	AuditFields `json:",inline"`
}
