package model

import "github.com/shappy0/mantritva/pkg/utils"

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	UserID   string  `json:"user_id" gorm:"uniqueIndex"`
	Mobile   string  `json:"mobile"`
	Email    string  `json:"email" gorm:"uniqueIndex"`
	Password *string `json:"password,omitempty"`
	// Role            string  `json:"role"`
	TokenVersion    int  `json:"token_version" gorm:"default:0"`
	AccountVerified bool `json:"account_verified" gorm:"default:false"`
	AuditFields
}

type Roles struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	AuditFields `json:",inline"`
}

type UserRole struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	UserID      uint            `json:"user_id" gorm:"uniqueIndex:idx_user_role"`
	RoleID      uint            `json:"role_id" gorm:"uniqueIndex:idx_user_role"`
	Role        utils.Roles     `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	Scope       utils.RoleScope `json:"scope" gorm:"not null"`
	ScopeID     uint            `json:"scope_id,omitempty" gorm:"uniqueIndex:idx_user_role;default:null"`
	AuditFields `json:",inline"`
}
