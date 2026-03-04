package model

import (
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key Result Area (KRA) represents a specific area of focus or responsibility within an organization that contributes to the achievement of its goals and objectives.
type KRA struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	KRAName     string `json:"kra_name"`
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
	TotalScore      float64            `json:"total_score" bson:"total_score"`
	MaxScore        uint               `json:"max_score" bson:"max_score"`
	Year            string             `json:"year" bson:"year"`
	Confidence      float64            `json:"confidence" bson:"confidence"`
	Developments    DevelopmentReport  `json:"developments" bson:"developments"`
	Grievances      GrievanceReport    `json:"grievances" bson:"grievances"`
	NewsPresses     NewsPressReport    `json:"news_presses" bson:"news_presses"`
	CALRecords      CALRecordsReport   `json:"cal_records" bson:"cal_records"`
	AuditFieldsBSON `bson:",inline"`
}

type MemberReportKRA struct {
	Developments DevelopmentReport `json:"developments" bson:"developments"`
	Grievances   GrievanceReport   `json:"grievances" bson:"grievances"`
	NewsPresses  NewsPressReport   `json:"news_presses" bson:"news_presses"`
	CALRecords   CALRecordsReport  `json:"cal_records" bson:"cal_records"`
}

type DevelopmentReport struct {
	Score         float64  `json:"score" bson:"score"`
	Weight        uint     `json:"weight" bson:"weight"`
	Total         uint     `json:"total" bson:"total"`
	Confidence    float64  `json:"confidence" bson:"confidence"`
	GlobalAverage float64  `json:"global_average" bson:"global_average"`
	Developments  []string `json:"developments" bson:"developments"` //Ids of developments, can be used to fetch details if needed
}

type GrievanceReport struct {
	Score         float64  `json:"score" bson:"score"`
	Weight        uint     `json:"weight" bson:"weight"`
	Total         uint     `json:"total" bson:"total"`
	Confidence    float64  `json:"confidence" bson:"confidence"`
	GlobalAverage float64  `json:"global_average" bson:"global_average"`
	Grievances    []string `json:"grievances" bson:"grievances"`
}

type NewsPressReport struct {
	Score         float64  `json:"score" bson:"score"`
	Weight        uint     `json:"weight" bson:"weight"`
	Total         uint     `json:"total" bson:"total"`
	Confidence    float64  `json:"confidence" bson:"confidence"`
	GlobalAverage float64  `json:"global_average" bson:"global_average"`
	NewsPresses   []string `json:"news_presses" bson:"news_presses"`
}

type CALRecordsReport struct {
	Score         float64  `json:"score" bson:"score"`
	Weight        uint     `json:"weight" bson:"weight"`
	Total         uint     `json:"total" bson:"total"`
	Confidence    float64  `json:"confidence" bson:"confidence"`
	GlobalAverage float64  `json:"global_average" bson:"global_average"`
	CALRecords    []string `json:"cal_records" bson:"cal_records"`
}
