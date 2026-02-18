package model

import (
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KSA struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	KSAName     string `json:"ksa_name"`
	Weight      uint   `json:"weight"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	AuditFields
}

type Questionnaire struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	FeedType utils.FeedType `json:"feed_type"`
	Category string         `json:"category"`
	Title    string         `json:"title"`
	VoteType utils.Vote     `json:"vote_type"`
	Weight   int            `json:"weight"`
	AuditFields
}

type MemberReport struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	MemberID        uint               `json:"member_id" bson:"member_id"`
	TotalScore      uint               `json:"total_score" bson:"total_score"`
	MaxScore        uint               `json:"max_score" bson:"max_score"`
	Year            string             `json:"year" bson:"year"`
	Developments    DevelopmentReport  `json:"developments" bson:"developments"`
	Grievances      GrievanceReport    `json:"grievances" bson:"grievances"`
	NewsPresses     NewsPressReport    `json:"news_presses" bson:"news_presses"`
	CALRecords      CALRecordsReport   `json:"cal_records" bson:"cal_records"`
	AuditFieldsBSON `bson:",inline"`
}

type MemberReportKSA struct {
	Developments DevelopmentReport `json:"developments" bson:"developments"`
	Grievances   GrievanceReport   `json:"grievances" bson:"grievances"`
	NewsPresses  NewsPressReport   `json:"news_presses" bson:"news_presses"`
	CALRecords   CALRecordsReport  `json:"cal_records" bson:"cal_records"`
}

type DevelopmentReport struct {
	Score        uint     `json:"score" bson:"score"`
	Weight       uint     `json:"weight" bson:"weight"`
	Total        uint     `json:"total" bson:"total"`
	Developments []string `json:"developments" bson:"developments"`
}

type GrievanceReport struct {
	Score      uint     `json:"score" bson:"score"`
	Weight     uint     `json:"weight" bson:"weight"`
	Total      uint     `json:"total" bson:"total"`
	Grievances []string `json:"grievances" bson:"grievances"`
}

type NewsPressReport struct {
	Score       uint     `json:"score" bson:"score"`
	Weight      uint     `json:"weight" bson:"weight"`
	Total       uint     `json:"total" bson:"total"`
	NewsPresses []string `json:"news_presses" bson:"news_presses"`
}

type CALRecordsReport struct {
	Score      uint     `json:"score" bson:"score"`
	Weight     uint     `json:"weight" bson:"weight"`
	Total      uint     `json:"total" bson:"total"`
	CALRecords []string `json:"cal_records" bson:"cal_records"`
}
