package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type LeadershipReq struct {
	RulerID       uint           `json:"ruler_id"`
	MemberID      uint           `json:"member_id"`
	PartyID       uint           `json:"party_id"`
	Position      string         `json:"position"`
	Role          string         `json:"role,omitempty"`
	Departments   string         `json:"departments,omitempty"`
	StateID       uint           `json:"state_id" required:"false"`
	LoksabhaID    uint           `json:"loksabha_id,omitempty"`
	VidhansabhaID uint           `json:"vidhansabha_id,omitempty"`
	RajyasabhaID  uint           `json:"rajyasabha_id,omitempty"`
	IsCentral     bool           `json:"is_central"`
	TermStart     utils.DateOnly `json:"term_start" required:"true"`
	TermEnd       utils.DateOnly `json:"term_end,omitempty"`
	IsActive      bool           `json:"is_active"`
}

type Leader struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Designation Position       `json:"designation"`
	Photo       string         `json:"photo"`
	State       string         `json:"state,omitempty"`
	Loksabha    string         `json:"loksabha,omitempty"`
	Vidhansabha string         `json:"vidhansabha,omitempty"`
	Rajyasabha  string         `json:"rajyasabha,omitempty"`
	TermStart   utils.DateOnly `json:"term_start"`
	TermEnd     utils.DateOnly `json:"term_end,omitempty"`
	Party       Party          `json:"party"`
}

type Position struct {
	Position    string  `json:"position"`
	Role        string  `json:"role,omitempty"`
	Departments string  `json:"departments,omitempty"`
	Region      *Region `json:"region,omitempty"`
}

type LeadershipRes struct {
	model.Leadership
	Member   Member   `json:"member"`
	Party    Party    `json:"party"`
	State    State    `json:"state,omitempty"`
	Position Position `json:"position,omitempty"`
	Region   Region   `json:"region,omitempty"`
}

type Region struct {
	State       string `json:"state,omitempty"`
	District    string `json:"district,omitempty"`
	Loksabha    string `json:"loksabha,omitempty"`
	Vidhansabha string `json:"vidhansabha,omitempty"`
	Rajyasabha  string `json:"rajyasabha,omitempty"`
	LocalBody   string `json:"local_body,omitempty"`
}

type LeadershipExitReq struct {
	LeadershipID uint           `json:"leadership_id"`
	TermEnd      utils.DateOnly `json:"term_end"`
}
