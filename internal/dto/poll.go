package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PollCreateReq struct {
	Question     string   `json:"question"`
	PollType     string   `json:"poll_type"`
	Options      []string `json:"options"`
	DurationDays int      `json:"duration_days"`
}
type PollVoteReq struct {
	OptionID int `json:"option_id" binding:"required"`
}
type PollOptionResult struct {
	ID         int     `json:"id"`
	Text       string  `json:"text"`
	Votes      int64   `json:"votes"`
	Percentage float64 `json:"percentage"`
}
type PollRes struct {
	ID               primitive.ObjectID `json:"id"`
	Question         string             `json:"question"`
	PollType         string             `json:"poll_type"`
	Options          []PollOptionResult `json:"options"`
	StartsAt         time.Time          `json:"starts_at"`
	EndsAt           time.Time          `json:"ends_at"`
	Status           string             `json:"status"`
	HasVoted         bool               `json:"has_voted"`
	SelectedOptionID *int               `json:"selected_option_id,omitempty"`
}

func PollOptions(p *model.Poll) []PollOptionResult {
	out := make([]PollOptionResult, len(p.Options))
	for i, o := range p.Options {
		out[i] = PollOptionResult{ID: o.ID, Text: o.Text}
	}
	return out
}
