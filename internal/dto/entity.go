package dto

import "github.com/shappy0/mantritva/internal/model"

type EntityReq struct {
	Name       string `json:"name" binding:"required"`
	Abbrv      string `json:"abbrv,omitempty"`
	IsCentral  *bool  `json:"is_central" binding:"required"`
	StateID    int    `json:"state_id,omitempty"`
	DistrictID int    `json:"district_id,omitempty"`
}

type EntityRes struct {
	model.Entity
	StateName    string `json:"state_name,omitempty"`
	DistrictName string `json:"district_name,omitempty"`
}

type Entity struct {
	Name         string `json:"name"`
	Abbrv        string `json:"abbrv,omitempty"`
	IsCentral    bool   `json:"is_central"`
	StateName    string `json:"state_name,omitempty"`
	DistrictName string `json:"district_name,omitempty"`
}
