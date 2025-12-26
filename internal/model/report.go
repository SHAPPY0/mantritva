package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type KSA struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	KSAName     string `json:"ksa_name"`
	Weight      uint   `json:"weight"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
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
	FIRs            FIRReport          `json:"firs" bson:"firs"`
	AuditFieldsBSON `bson:",inline"`
}

type MemberReportKSA struct {
	Developments DevelopmentReport `json:"developments" bson:"developments"`
	Grievances   GrievanceReport   `json:"grievances" bson:"grievances"`
	NewsPresses  NewsPressReport   `json:"news_presses" bson:"news_presses"`
	FIRs         FIRReport         `json:"firs" bson:"firs"`
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

type FIRReport struct {
	Score  uint     `json:"score" bson:"score"`
	Weight uint     `json:"weight" bson:"weight"`
	Total  uint     `json:"total" bson:"total"`
	FIRs   []string `json:"firs" bson:"firs"`
}
