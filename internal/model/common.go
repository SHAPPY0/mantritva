package model

import (
	"time"

	"gorm.io/gorm"
)

type AuditFields struct {
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedBy string         `json:"updated_by"`
	IsDeleted bool           `json:"is_deleted" gorm:"default:false"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type AuditFieldsBSON struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy string    `json:"updated_by" bson:"updated_by"`
	IsDeleted bool      `json:"is_deleted" bson:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`
}

func InsertBSONAuditFields(createdBy string) AuditFieldsBSON {
	a := AuditFieldsBSON{
		CreatedAt: time.Now(),
		CreatedBy: "SYSTEM",
		IsDeleted: false,
	}
	if createdBy != "" {
		a.CreatedBy = createdBy
	}
	return a
}
