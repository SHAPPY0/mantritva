package dto

import (
	"mime/multipart"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrievanceReq struct {
	Type                 string                  `json:"type" form:"type"`
	SubType              string                  `json:"subtype" form:"subtype"`
	Jurisdiction         string                  `json:"jurisdiction" form:"jurisdiction"`
	CenterSubLevel       string                  `json:"center_sublevel" form:"center_sublevel"`
	CabinetMinisterID    uint                    `json:"cabinet_minister_id" form:"cabinet_minister_id" binding:"omitempty"`
	MemberOfParliamentID uint                    `json:"member_of_parliament_id" form:"member_of_parliament_id" binding:"omitempty"`
	EntityID             uint                    `json:"entity_id" form:"entity_id" binding:"omitempty"`
	StateID              uint                    `json:"state_id" form:"state_id" binding:"omitempty"`
	DistrictID           uint                    `json:"district_id" form:"district_id" binding:"omitempty"`
	StateLeaderID        uint                    `json:"state_leader_id" form:"state_leader_id" binding:"omitempty"`
	Title                string                  `json:"title" form:"title"`
	Description          string                  `json:"description" form:"description"`
	HideIdentity         bool                    `json:"hide_identity" form:"hide_identity"`
	DateOfIncident       utils.DateOnly          `json:"date_of_incident" form:"date_of_incident"`
	Location             string                  `json:"location" form:"location"`
	Impact               string                  `json:"impact" form:"impact"`
	Severity             utils.GrievanceSeverity `json:"severity" form:"severity"`
	// Evidences            []EvidenceReq      `json:"evidences" form:"evidences"`
	SLA         *utils.DateOnly       `json:"sla" form:"sla"`
	FinalStatus utils.GrievanceStatus `json:"final_status" form:"final_status"`
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
	CabinetMinister    *Member          `json:"cabinet_minister,omitempty"`
	MemberOfParliament *Member          `json:"member_of_parliament,omitempty"`
	StateLeader        *Member          `json:"state_leader,omitempty"`
	State              *model.State     `json:"state,omitempty"`
	District           string           `json:"district,omitempty"`
	Entity             *Entity          `json:"entity,omitempty"`
	User               UserClaim        `json:"user,omitempty"`
	Reaction           utils.Vote       `json:"reaction"`
	ReactionLocked     bool             `json:"reaction_locked"`
	Questionnaire      *[]Questionnaire `json:"questionnaire,omitempty"`
}

type GrievanceStatus struct {
	ID     string                `json:"id"`
	Date   utils.DateOnly        `json:"date"`
	Status utils.GrievanceStatus `json:"status"`
	Remark string                `json:"remark"`
	Source string                `json:"source"`
}

type GrievanceShort struct {
	ID             primitive.ObjectID    `json:"id"`
	GrievanceID    string                `json:"grievance_id,omitempty"`
	Type           string                `json:"type,omitempty"`
	SubType        string                `json:"subtype,omitempty"`
	Jurisdiction   string                `json:"jurisdiction,omitempty"`
	CenterSubLevel string                `json:"center_sublevel,omitempty"`
	Title          string                `json:"title,omitempty"`
	Description    string                `json:"description"`
	DateOfIncident utils.DateOnly        `json:"date_of_incident,omitempty"`
	Location       string                `json:"location,omitempty"`
	Impact         string                `json:"impact,omitempty"`
	Severity       string                `json:"severity"`
	FinalStatus    utils.GrievanceStatus `json:"final_status,omitempty"`
	CreatedAt      time.Time             `json:"created_at,omitempty"`
	Questionnaire  *[]Questionnaire      `json:"questionnaire,omitempty"`
	Reaction       utils.Vote            `json:"reaction"`
	ReactionLocked bool                  `json:"reaction_locked"`
	UpvoteCount    uint                  `json:"upvote_count"`
	DownvoteCount  uint                  `json:"downvote_count"`
}
