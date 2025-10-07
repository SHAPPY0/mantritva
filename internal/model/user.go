package model

import "github.com/shappy0/mantritva/pkg/utils"

type User struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	Name     string         `json:"name"`
	UserID   string         `json:"user_id" gorm:"uniqueIndex"`
	Mobile   string         `json:"mobile"`
	Email    string         `json:"email" gorm:"uniqueIndex"`
	Password string         `json:"password"`
	Role     utils.UserRole `json:"role"`
	AuditFields
}
