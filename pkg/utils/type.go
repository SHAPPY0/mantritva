package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DateOnly struct {
	*time.Time
}
type MongoID primitive.ObjectID
type StatementType string
type Subject string
type StatementTone string
type PostStatus string
type UserRole string
type FeedType string
type PostReactionType string

const (
	ADMIN     UserRole = "ADMIN"
	MODERATOR UserRole = "MODERATOR"
	USER      UserRole = "USER"
)

const (
	PolicyAnnouncement StatementType = "PolicyAnnouncement"
	Reaction           StatementType = "Reaction"
	Criticism          StatementType = "Criticism"
	Defense            StatementType = "Defense"
	ElectionPromise    StatementType = "ElectionPromise"
	OtherStatement     StatementType = "Other"
)

const (
	Economy          Subject = "Economy"
	ForeignPolicy    Subject = "ForeignPolicy"
	SocialIssue      Subject = "SocialIssue"
	ElectionCampaign Subject = "ElectionCampaign"
	OppositionAttack Subject = "OppositionAttack"
	OtherSubject     Subject = "Other"
)

const (
	Supportive StatementTone = "Supportive"
	Critical   StatementTone = "Critical"
	Aggressive StatementTone = "Aggressive"
	Diplomatic StatementTone = "Diplomatic"
	Emotional  StatementTone = "Emotional"
	Other      StatementTone = "Other"
)

const (
	Draft     PostStatus = "Draft"
	Pending   PostStatus = "Pending"
	Published PostStatus = "Published"
	Archived  PostStatus = "Archived"
)

const (
	NEWSPRESS   FeedType = "NEWSPRESS"
	DEVELOPMENT FeedType = "DEVELOPMENT"
)

const (
	UPVOTE   PostReactionType = "UPVOTE"
	DOWNVOTE PostReactionType = "DOWNVOTE"
)

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = nil
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = &t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time == nil || d.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

func (d DateOnly) Value() (driver.Value, error) {
	if d.Time == nil {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}

func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok && strings.TrimSpace(s) == "" {
		return nil
	}
	if t, ok := value.(time.Time); ok {
		d.Time = &t
		return nil
	}
	return fmt.Errorf("cannot convert %v to DateOnly", value)
}

func (d *DateOnly) IsNull() bool {
	return d.Time == nil || d.Time.IsZero()
}

// Method to convert CustomObjectID to string
func (id MongoID) String() string {
	return primitive.ObjectID(id).Hex()
}

// Method to check if CustomObjectID is valid
func (id MongoID) IsValid() bool {
	return primitive.ObjectID(id).IsZero()
}
