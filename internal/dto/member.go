package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type MemberReq struct {
	Name             string              `json:"name" form:"name" binding:"required"`
	FullName         string              `json:"full_name" form:"full_name" binding:"required"`
	DOB              utils.DateOnly      `json:"dob" form:"dob" gorm:"type:date"`
	Gender           string              `json:"gender" form:"gender" binding:"required"`
	Qualifications   string              `json:"qualifications" form:"qualifications"`
	MaritalStatus    string              `json:"marital_status" form:"marital_status"`
	FinancialAssets  string              `json:"financial_assets" form:"financial_assets"`
	SocialMediaLinks map[string]string   `json:"social_media_links" form:"social_media_links" binding:"required"`
	CriminalRecords  []CriminalRecordReq `json:"criminal_records,omitempty" form:"criminal_records,omitempty"`
}

type CriminalRecordReq struct {
	Offense        string         `form:"offense"`
	Status         string         `form:"status"`
	ConvictionDate utils.DateOnly `form:"conviction_date" gorm:"type:date"`
	Sentence       string         `form:"sentence"`
}

type CareerReq struct {
	PartyID       uint           `form:"party_id" binding:"required"`
	Position      string         `form:"position,omitempty" binding:"required"`
	TermStartDate utils.DateOnly `form:"term_start_date" gorm:"type:date" binding:"required"`
	TermEndDate   utils.DateOnly `form:"term_end_date,omitempty" gorm:"type:date"`
}

type CareerRes struct {
	Party     Party          `json:"party"`
	TermStart utils.DateOnly `json:"term_start"`
	TermEnd   utils.DateOnly `json:"term_end,omitempty"`
	Position  Position       `json:"position"`
	Region    Region         `json:"region,omitempty"`
}

type MemberParty struct {
	Name        string         `json:"name,omitempty"`
	Abbrv       string         `json:"abbrv,omitempty"`
	Symbol      string         `json:"symbol,omitempty"`
	Position    string         `json:"position,omitempty"`
	JoiningDate utils.DateOnly `json:"joining_date,omitempty"`
	ExitDate    utils.DateOnly `json:"exit_date,omitempty"`
}

type MemberRes struct {
	model.Member
	Position    string       `json:"position,omitempty"`
	Departments string       `json:"departments,omitempty"`
	Party       MemberParty  `json:"party,omitempty"`
	Score       *MemberScore `json:"score,omitempty"`
}

type Member struct {
	Name     string      `json:"name"`
	Position string      `json:"position"`
	Photo    string      `json:"photo"`
	Party    MemberParty `json:"party,omitempty"`
}
