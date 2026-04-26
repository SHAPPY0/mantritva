package dto

import "github.com/shappy0/mantritva/pkg/utils"

type ULBReq struct {
	Name       string        `json:"name"`
	Code       string        `json:"code,omitempty"`
	ShortName  string        `json:"short_name,omitempty"` //Can store city name
	DistrictID uint          `json:"district_id"`
	BodyType   utils.ULBType `json:"type"`
	TotalWards uint          `json:"total_wards"`
}
