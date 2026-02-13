package dto

import (
	"mime/multipart"

	"github.com/shappy0/mantritva/pkg/utils"
)

type DevelopmentWelfareReq struct {
	Category               string                  `json:"category" form:"category"`
	Type                   string                  `json:"type" form:"type"`
	PartyID                uint                    `json:"party_id" form:"party_id"`
	MemberID               uint                    `json:"member_id" form:"member_id"`
	Name                   string                  `json:"name" form:"name"`
	Description            string                  `json:"description" form:"description"`
	StateID                uint                    `json:"state_id" form:"state_id"`
	DistrictID             uint                    `json:"district_id" form:"district_id"`
	VidhanSabhaID          uint                    `json:"vidhansabha_id" form:"vidhansabha_id"`
	LokSabhaID             uint                    `json:"loksabha_id" form:"loksabha_id"`
	Location               string                  `json:"location" form:"location"`
	AnnouncementDate       utils.DateOnly          `json:"announcement_date" form:"announcement_date"`
	ExpectedCompletionDate utils.DateOnly          `json:"expected_completion_date" form:"expected_completion_date"`
	ActualCompletionDate   utils.DateOnly          `json:"actual_completion_date" form:"actual_completion_date"`
	Criteria               string                  `json:"criteria" form:"criteria"`
	BeneficiariesNum       uint                    `json:"beneficiaries_num" form:"beneficiaries_num"`
	Scope                  string                  `json:"scope" form:"scope"`
	Status                 utils.DevelopmentStatus `json:"status" form:"status"`
	Documents              []Document              `json:"documents"`
	CreatedBy              string                  `json:"created_by"`
}

type Document struct {
	File     *multipart.FileHeader `json:"file"`
	FileName string                `json:"filename"`
	Remark   string                `json:"remark"`
}

type DevelopmentWelfareShort struct {
	ID               string                  `json:"_id,omitempty"`
	Category         string                  `json:"category"`
	Type             string                  `json:"type"`
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	Member           Member                  `json:"member,omitempty"`
	Party            Party                   `json:"party,omitempty"`
	AnnouncementDate utils.DateOnly          `json:"announcement_date"`
	Status           utils.DevelopmentStatus `json:"status"`
	Reaction         utils.Vote              `json:"reaction"`
	UpvoteCount      uint                    `json:"upvote_count"`
	DownvoteCount    uint                    `json:"downvote_count"`
	CreatedAt        utils.DateOnly          `json:"created_at"`
}
