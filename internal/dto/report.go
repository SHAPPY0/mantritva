package dto

import (
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KSAWeights struct {
	Developments uint
	Grievances   uint
	FIRs         uint
	NewsPresses  uint
}

type MemberReport struct {
	MemberID     uint              `json:"member_id"`
	TotalScore   uint              `json:"total_score"`
	Tenure       string            `json:"tenure"`
	Developments DevelopmentReport `json:"developments"`
	Grievances   GrievanceReport   `json:"grievances"`
	NewsPresses  NewsPressesReport `json:"news_presses"`
}

type DevelopmentReport struct {
	Score uint `json:"score"`
}

type GrievanceReport struct {
	Score      uint             `json:"score"`
	Grievances []GrievanceShort `json:"grievances"`
}

type NewsPressesReport struct {
	Score       uint             `json:"score"`
	NewsPresses []NewsPressShort `json:"news_presses"`
}

type MemberScore struct {
	Score     uint   `json:"score"`
	MaxScore  uint   `json:"max_score"`
	Year      string `json:"year"`
	UpdatedAt string `json:"updated_at"`
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
