package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LocalBodyRes struct {
	model.ULB     `json:",inline,omitempty"`
	LocalBodyType utils.LocalBodyType `json:"local_body_type"`
	DisplayType   string              `json:"display_type"`
}

type ULBReq struct {
	Name       string        `json:"name"`
	Code       string        `json:"code,omitempty"`
	Area       string        `json:"area,omitempty"` //Can store city name
	DistrictID uint          `json:"district_id"`
	BodyType   utils.ULBType `json:"type"`
	TotalWards uint          `json:"total_wards"`
}

type ULBShort struct {
	ID       uint          `json:"id"`
	Name     string        `json:"name"`
	BodyType utils.ULBType `json:"type"`
}
