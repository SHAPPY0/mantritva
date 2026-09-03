package model

import "github.com/shappy0/mantritva/pkg/utils"

type LocalBody struct {
	ID          uint                `json:"id,omitempty" gorm:"primaryKey"`
	Name        string              `json:"name"`
	Area        string              `json:"area,omitempty"`
	Code        string              `json:"code,omitempty"`
	DistrictID  uint                `json:"district_id"`
	BodyType    utils.LocalBodyType `json:"type"`
	TotalWards  uint                `json:"total_wards" gorm:"default:0"`
	AuditFields `json:",inline"`
}

type Ward struct {
	ID          uint   `json:"id,omitempty" gorm:"primaryKey"`
	Name        string `json:"name"`
	Number      uint   `json:"number,omitempty"`
	LocalBodyID uint   `json:"local_body_id,omitempty"`
	AuditFields `json:",inline"`
}
