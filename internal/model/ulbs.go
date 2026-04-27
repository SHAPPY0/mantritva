package model

import "github.com/shappy0/mantritva/pkg/utils"

type ULB struct {
	ID          uint          `json:"id,omitempty" gorm:"primaryKey"`
	Name        string        `json:"name"`
	Area        string        `json:"area,omitempty"`
	Code        string        `json:"code,omitempty"`
	DistrictID  uint          `json:"district_id"`
	BodyType    utils.ULBType `json:"type"`
	TotalWards  uint          `json:"total_wards" gorm:"default:0"`
	AuditFields `json:",inline"`
}

type Ward struct {
	ID          uint   `json:"id,omitempty" gorm:"primaryKey"`
	Name        string `json:"name"`
	Number      uint   `json:"number,omitempty"`
	ULBID       uint   `json:"ulb_id,omitempty"`
	AuditFields `json:",inline"`
}
