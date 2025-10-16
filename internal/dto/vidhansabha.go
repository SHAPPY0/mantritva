package dto

type VidhanReq struct {
	Name       string `json:"name" binding:"required"`
	StateID    int    `json:"state_id" binding:"required"`
	DistrictID int    `json:"district_id" binding:"required"`
	ACNumber   int    `json:"ac_number,omitempty"`
}

type VidhanSabhaRes struct {
	Name     string `json:"name"`
	District string `json:"district"`
	State    string `json:"state"`
}
