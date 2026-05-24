package dto

import (
	"github.com/shappy0/mantritva/pkg/utils"
)

type Enquiry struct {
	Name        string `form:"name"`
	Email       string `form:"email"`
	Mobile      string `form:"mobile"`
	EnquiryType string `form:"enquiry_type"`
	Subject     string `form:"subject"`
	Description string `form:"description"`
}

type BulkLeader struct {
	Member MemberReq `json:"member"`
	Party  struct {
		Setup       bool           `json:"setup"`
		Name        string         `json:"name"`
		Position    string         `json:"position"`
		JoiningDate utils.DateOnly `json:"joining_date"`
		ExitDate    utils.DateOnly `json:"exit_date,omitempty"`
	}
	Leadership struct {
		Setup       bool           `json:"setup"`
		State       string         `json:"state"`
		Position    string         `json:"position"`
		Role        string         `json:"role"`
		Departments string         `json:"departments"`
		TermStart   utils.DateOnly `json:"term_start"`
		TermEnd     utils.DateOnly `json:"term_end,omitempty"`
		Loksabha    string         `json:"loksabha,omitempty"`
		Vidhansabha string         `json:"vidhansabha,omitempty"`
		Rajyasabha  string         `json:"rajyasabha,omitempty"`
		IsCentral   bool           `json:"is_central"`
		IsActive    bool           `json:"is_active"`
	}
}
