package dto

type CountryReq struct {
	Name    string `json:"name" binding:"required"`
	Abbrv   string `json:"abbrv" binding:"required"`
	Icon    string `json:"icon"`
	Capital string `json:"capital" binding:"required"`
}
