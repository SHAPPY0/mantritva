package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LeadershipReq struct {
	RulerID     uint           `json:"ruler_id"`
	MemberID    uint           `json:"member_id"`
	PartyID     uint           `json:"party_id"`
	Position    string         `json:"position"`
	Role        string         `json:"role,omitempty"`
	Departments string         `json:"departments,omitempty"`
	StateID     uint           `json:"state_id" required:"false"`
	IsCentral   bool           `json:"is_central"`
	TermStart   utils.DateOnly `json:"term_start" required:"true"`
	TermEnd     utils.DateOnly `json:"term_end,omitempty"`
	IsActive    bool           `json:"is_active"`
}

type Leader struct {
	Name        string         `json:"name"`
	Designation Position       `json:"designation"`
	Photo       string         `json:"photo"`
	TermStart   utils.DateOnly `json:"term_start"`
	TermEnd     utils.DateOnly `json:"term_end,omitempty"`
}

type Position struct {
	Position    string `json:"position"`
	Role        string `json:"role"`
	Departments string `json:"departments,omitempty"`
}

type LeadershipRes struct {
	model.Leadership
	Member Member `json:"member"`
	Party  Party  `json:"party"`
	State  State  `json:"state"`
}
