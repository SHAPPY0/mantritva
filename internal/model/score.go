package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Score struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	MemberID        uint               `json:"member_id" bson:"member_id"`
	Year            string             `json:"year" bson:"year"`
	TotalScore      uint               `json:"total_score" bson:"total_score"`
	MaxScore        uint               `json:"max_score" bson:"max_score"`
	Confidence      uint               `json:"confidence" bson:"confidence"`
	Developments    DevelopmentScore   `json:"developments_score" bson:"developments_score"`
	Grievances      GrievanceScore     `json:"grievances_score" bson:"grievances_score"`
	NewsPresses     NewsPressScore     `json:"news_presses_score" bson:"news_presses_score"`
	CALRecords      CALRecordScore     `json:"cal_records_score" bson:"cal_records_score"`
	Version         uint               `json:"version" bson:"version"`
	AuditFieldsBSON `bson:",inline"`
}

type DevelopmentScore struct {
	Score          uint      `json:"score" bson:"score"`
	Weight         uint      `json:"weight" bson:"weight"`
	TotalCount     uint      `json:"total_count" bson:"total_count"`
	Confidence     uint      `json:"confidence" bson:"confidence"`
	GlobalAverage  uint      `json:"global_average" bson:"global_average"`
	DevelopmentIds []string  `json:"development_ids" bson:"development_ids"` //Ids of developments, can be used to fetch details if needed
	LastActivityAt time.Time `json:"last_activity_at" bson:"last_activity_at"`
}

type GrievanceScore struct {
	Score          uint      `json:"score" bson:"score"`
	Weight         uint      `json:"weight" bson:"weight"`
	TotalCount     uint      `json:"total_count" bson:"total_count"`
	Confidence     uint      `json:"confidence" bson:"confidence"`
	GlobalAverage  uint      `json:"global_average" bson:"global_average"`
	GrievanceIds   []string  `json:"grievance_ids" bson:"grievance_ids"` //Ids of grievances, can be used to fetch details if needed
	LastActivityAt time.Time `json:"last_activity_at" bson:"last_activity_at"`
}

type NewsPressScore struct {
	Score          uint      `json:"score" bson:"score"`
	Weight         uint      `json:"weight" bson:"weight"`
	TotalCount     uint      `json:"total_count" bson:"total_count"`
	Confidence     uint      `json:"confidence" bson:"confidence"`
	GlobalAverage  uint      `json:"global_average" bson:"global_average"`
	NewsPressIds   []string  `json:"news_press_ids" bson:"news_press_ids"` //Ids of news presses, can be used to fetch details if needed
	LastActivityAt time.Time `json:"last_activity_at" bson:"last_activity_at"`
}

type CALRecordScore struct {
	Score          uint      `json:"score" bson:"score"`
	Weight         uint      `json:"weight" bson:"weight"`
	TotalCount     uint      `json:"total_count" bson:"total_count"`
	Confidence     uint      `json:"confidence" bson:"confidence"`
	GlobalAverage  uint      `json:"global_average" bson:"global_average"`
	CALRecordIds   []string  `json:"cal_record_ids" bson:"cal_record_ids"` //Ids of cal records, can be used to fetch details if needed
	LastActivityAt time.Time `json:"last_activity_at" bson:"last_activity_at"`
}
