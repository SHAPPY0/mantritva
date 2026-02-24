package model

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	UserID       string `json:"user_id" gorm:"uniqueIndex"`
	Mobile       string `json:"mobile"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version" gorm:"default:0"`
	AuditFields
}
