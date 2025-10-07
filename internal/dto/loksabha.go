package dto

type LoksabhaReq struct {
	Name       string `json:"name" binding:"required"`
	StateID    int    `json:"state_id" binding:"required"`
	DistrictID int    `json:"district_id,omitempty"`
	PCNumber   int    `json:"pc_number" binding:"required"`
}
