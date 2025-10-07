package dto

type VidhanReq struct {
	Name       string `json:"name" binding:"required"`
	LoksabhaID int    `json:"loksabha_id" binding:"required"`
	ACNumber   int    `json:"ac_number" binding:"required"`
}
