package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type MemberReq struct {
	Name             string              `form:"name" binding:"required"`
	FullName         string              `form:"full_name" binding:"required"`
	DOB              utils.DateOnly      `form:"dob" gorm:"type:date"`
	Gender           string              `form:"gender" binding:"required"`
	Qualifications   string              `form:"qualifications"`
	MaritalStatus    string              `form:"marital_status"`
	FinancialAssets  string              `form:"financial_assets"`
	SocialMediaLinks map[string]string   `form:"social_media_links" binding:"required"`
	CriminalRecords  []CriminalRecordReq `form:"criminal_records,omitempty"`
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

type MemberParty struct {
	Name          string         `json:"name,omitempty"`
	Abbrv         string         `json:"abbrv,omitempty"`
	Symbol        string         `json:"symbol,omitempty"`
	Position      string         `json:"position,omitempty"`
	TermStartDate utils.DateOnly `json:"term_start_date,omitempty"`
	TermEndDate   utils.DateOnly `json:"term_end_date,omitempty"`
}

type MemberRes struct {
	model.Member
	Position    string      `json:"position,omitempty"`
	Departments string      `json:"departments,omitempty"`
	Party       MemberParty `json:"party,omitempty"`
}

type Member struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Photo    string `json:"photo"`
}
