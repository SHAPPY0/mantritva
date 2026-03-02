package model

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Grievance struct {
	ID                   primitive.ObjectID    `json:"_id" bson:"_id,omitempty"`
	GrievanceID          string                `json:"grievance_id" bson:"grievance_id"`
	Type                 string                `json:"type" bson:"type"`
	SubType              string                `json:"subtype" bson:"subtype"`
	Level                string                `json:"level" bson:"level"`
	CenterSubLevel       string                `json:"center_sublevel" bson:"center_sublevel"`
	CabinetMinisterID    uint                  `json:"cabinet_minister_id" bson:"cabinet_minister_id"`
	MemberOfParliamentID uint                  `json:"member_of_parliament_id" bson:"member_of_parliament_id"`
	StateID              uint                  `json:"state_id" bson:"state_id"`
	DistrictID           uint                  `json:"district_id" bson:"district_id"`
	StateLeaderID        uint                  `json:"state_leader_id" bson:"state_leader_id"`
	Title                string                `json:"title" bson:"title"`
	Description          string                `json:"description" bson:"description"`
	HideIdentity         bool                  `json:"hide_identity" bson:"hide_identity"`
	DateOfIncident       time.Time             `json:"date_of_incident" bson:"date_of_incident"`
	Location             string                `json:"location" bson:"location"`
	Impact               string                `json:"impact" bson:"impact"`
	Evidences            []Document            `json:"evidences" bson:"evidences"`
	FinalStatus          utils.GrievanceStatus `json:"final_status" bson:"final_status"`
	StatusRevisions      []StatusRevisions     `json:"status_revisions" bson:"status_revisions"`
	AuditFieldsBSON      `bson:",inline"`
}

type Document struct {
	Url      string `json:"url" bson:"url"`
	FileName string `json:"filename" bson:"filename"`
	FileType string `json:"filetype" bson:"filetype"`
	Remark   string `json:"remark" bson:"remark"`
}

type StatusRevisions struct {
	Date   time.Time `json:"date" bson:"date"`
	Status string    `json:"status" bson:"status"`
	Remark string    `json:"remark" bson:"remark"`
	Source string    `json:"source" bson:"source"`
}
