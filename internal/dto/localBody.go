package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LocalBodyRes struct {
	model.LocalBody `json:",inline,omitempty"`
	DisplayType     string `json:"display_type"`
}

type LocalBodyReq struct {
	Name       string              `json:"name"`
	Code       string              `json:"code,omitempty"`
	Area       string              `json:"area,omitempty"`
	DistrictID uint                `json:"district_id"`
	BodyType   utils.LocalBodyType `json:"type"`
	TotalWards uint                `json:"total_wards"`
}

type LocalBodyShort struct {
	ID       uint                `json:"id"`
	Name     string              `json:"name"`
	BodyType utils.LocalBodyType `json:"body_type"`
	District string              `json:"district"`
}
