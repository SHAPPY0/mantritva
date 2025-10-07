package model

import (
	"time"

	"gorm.io/gorm"
)

type AuditFields struct {
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	IsDeleted bool           `json:"is_deleted" gorm:"default:false"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type AuditFieldsBSON struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	IsDeleted bool      `json:"is_deleted" bson:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`
}

func InsertBSONAuditFields() AuditFieldsBSON {
	return AuditFieldsBSON{
		CreatedAt: time.Now(),
		CreatedBy: "SYSTEM",
		IsDeleted: false,
	}
}
