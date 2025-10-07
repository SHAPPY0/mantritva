package dto

type StateReq struct {
	Name      string `json:"name" binding:"required"`
	CountryID int    `json:"country_id" binding:"required"`
	Abbrv     string `json:"abbrv" binding:"required"`
	Capital   string `json:"capital" binding:"required"`
	Structure string `json:"structure" binding:"required"`
}

type State struct {
	Name  string `json:"name"`
	Abbrv string `json:"abbrv"`
}
