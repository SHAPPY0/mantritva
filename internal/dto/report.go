package dto

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KRAWeights struct {
	Developments uint
	Grievances   uint
	CALRecords   uint
	NewsPresses  uint
}

type MemberReport struct {
	MemberID   uint      `json:"member_id"`
	TotalScore float64   `json:"total_score"`
	Confidence float64   `json:"confidence"`
	Year       string    `json:"year"`
	KRAS       KRAReport `json:"kras"`
	CreatedAt  time.Time `json:"created_at"`
}

type KRAReport struct {
	Developments DevelopmentReport `json:"developments"`
	Grievances   GrievanceReport   `json:"grievances"`
	NewsPresses  NewsPressesReport `json:"news_pressses"`
	CALRecords   CALReport         `json:"cal_records"`
}

type DevelopmentReport struct {
	Name         string   `json:"name"`
	Score        float64  `json:"score"`
	Confidence   float64  `json:"confidence"`
	Weight       uint     `json:"weight"`
	Total        uint     `json:"total"`
	Developments []string `json:"developments"`
}

type GrievanceReport struct {
	Name       string   `json:"name"`
	Score      float64  `json:"score"`
	Confidence float64  `json:"confidence"`
	Weight     uint     `json:"weight"`
	Total      uint     `json:"total"`
	Grievances []string `json:"grievances"`
}

type NewsPressesReport struct {
	Name        string   `json:"name"`
	Score       float64  `json:"score"`
	Confidence  float64  `json:"confidence"`
	Weight      uint     `json:"weight"`
	Total       uint     `json:"total"`
	NewsPresses []string `json:"news_presses"`
}

type CALReport struct {
	Name       string   `json:"name"`
	Score      float64  `json:"score"`
	Confidence float64  `json:"confidence"`
	Weight     uint     `json:"weight"`
	Total      uint     `json:"total"`
	CALRecords []string `json:"cal_records"`
}

type MemberScore struct {
	Score      float64 `json:"score"`
	MaxScore   uint    `json:"max_score"`
	Confidence float64 `json:"confidence"`
	Year       *string `json:"year,omitempty"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}

type Questionnaire struct {
	ID       uint            `json:"id"`
	Category string          `json:"category,omitempty"`
	FeedType *utils.FeedType `json:"feed_type,omitempty"`
	Title    string          `json:"title"`
	VoteType *utils.Vote     `json:"vote_type"`
	Weight   *int            `json:"weight,omitempty"`
}

type FeedReactionItem struct {
	FeedType utils.FeedType
	Reaction string
	UserID   int
	PostID   primitive.ObjectID
}
