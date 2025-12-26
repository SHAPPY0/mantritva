package dto

import (
	"mime/multipart"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type GrievanceReq struct {
	Type                 string         `json:"type" form:"type"`
	SubType              string         `json:"subtype" form:"subtype"`
	Level                string         `json:"level" form:"level"`
	CenterSubLevel       string         `json:"center_sublevel" form:"center_sublevel"`
	CabinetMinisterID    uint           `json:"cabinet_minister_id" form:"cabinet_minister_id" binding:"omitempty"`
	MemberOfParliamentID uint           `json:"member_of_parliament_id" form:"member_of_parliament_id" binding:"omitempty"`
	StateID              uint           `json:"state_id" form:"state_id" binding:"omitempty"`
	DistrictID           uint           `json:"district_id" form:"district_id" binding:"omitempty"`
	StateLeaderID        uint           `json:"state_leader_id" form:"state_leader_id" binding:"omitempty"`
	Title                string         `json:"title" form:"title"`
	Description          string         `json:"description" form:"description"`
	HideIdentity         bool           `json:"hide_identity" form:"hide_identity"`
	DateOfIncident       utils.DateOnly `json:"date_of_incident" form:"date_of_incident"`
	Location             string         `json:"location" form:"location"`
	Impact               string         `json:"impact" form:"impact"`
	// Evidences            []EvidenceReq      `json:"evidences" form:"evidences"`
	FinalStatus utils.GrievanceStatus `json:"FinalStatus" form:"finalStatus"`
	// Updates              []GrievanceUpdates `json:"updates" form:"updates"`
	CreatedBy string `json:"created_by"`
}

type Grievance struct {
	GrievanceReq
	Evidences []EvidenceReq `json:"evidences"`
}

type EvidenceReq struct {
	FileName string                `json:"filename" form:"filename"`
	Comment  string                `json:"comment" form:"comment"`
	File     *multipart.FileHeader `json:"file" form:"file"`
}

type GrievanceRes struct {
	model.Grievance
	CabinetMinister    *Member      `json:"cabinet_minister,omitempty"`
	MemberOfParliament *Member      `json:"member_of_parliament,omitempty"`
	StateLeader        *Member      `json:"state_leader,omitempty"`
	State              *model.State `json:"state,omitempty"`
	District           string       `json:"district,omitempty"`
	User               UserClaim    `json:"user,omitempty"`
}

type GrievanceStatus struct {
	ID     string                `json:"id"`
	Date   utils.DateOnly        `json:"date"`
	Status utils.GrievanceStatus `json:"status"`
	Remark string                `json:"remark"`
	Link   string                `json:"link"`
}

type GrievanceShort struct {
	GrievanceID    string                `json:"grievance_id,omitempty"`
	Type           string                `json:"type,omitempty"`
	SubType        string                `json:"subtype,omitempty"`
	Level          string                `json:"level,omitempty"`
	CenterSubLevel string                `json:"center_sublevel,omitempty"`
	Title          string                `json:"title,omitempty"`
	DateOfIncident utils.DateOnly        `json:"date_of_incident,omitempty"`
	Location       string                `json:"location,omitempty"`
	Impact         string                `json:"impact,omitempty"`
	FinalStatus    utils.GrievanceStatus `json:"final_status,omitempty"`
	CreatedAt      time.Time             `json:"created_at,omitempty"`
}
