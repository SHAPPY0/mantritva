package dto

type DistrictReq struct {
	Name    string `json:"name" binding:"required"`
	Abbrv   string `json:"abbrv" binding:"required"`
	StateID int    `json:"state_id" binding:"required"`
}
