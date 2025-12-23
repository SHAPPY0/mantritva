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
type Position string
type GrievanceStatus string

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
	GRIEVANCE   FeedType = "GRIEVANCE"
)

const (
	UPVOTE   PostReactionType = "UPVOTE"
	DOWNVOTE PostReactionType = "DOWNVOTE"
)

const (
	PRIMEMINISTER       Position = "Prime Minister"
	PRESIDENT           Position = "President"
	VICEPRESIDENT       Position = "Vice President"
	CABINETMINISTER     Position = "Cabinet Minister"
	CHIEFMINISTER       Position = "Chief Minister"
	DEPUTYCHIEFMINISTER Position = "Deputy Chief Minister"
	OPPOSITIONLEADER    Position = "Opposition Leader"
	MP                  Position = "Member Of Parliament"
	MOS                 Position = "Member Of State"
	MLA                 Position = "Member Of Legislative Assembly"
	RAJYASABHAMEMEBER   Position = "Rajyasabha Member"
)

const (
	UNDERREVIEW GrievanceStatus = "UNDER_REVIEW"
	OPEN        GrievanceStatus = "OPEN"
	REJECTED    GrievanceStatus = "REJECTED"
	INPROGRESS  GrievanceStatus = "IN_PROGRESS"
	CLOSED      GrievanceStatus = "CLOSED"
	FIXED       GrievanceStatus = "FIXED"
	COMPLETED   GrievanceStatus = "COMPLETED"
	STAYED      GrievanceStatus = "STAYED"
	NOUPDATE    GrievanceStatus = "NO_UPDATE"
)

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = nil
		return nil
	}
	var err error
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		var t time.Time
		t, err = time.Parse(layout, s)
		if err == nil {
			d.Time = &t
			return nil
		}
	}
	return err
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
	return d == nil || d.Time == nil || d.Time.IsZero()
}

// Method to convert CustomObjectID to string
func (id MongoID) String() string {
	return primitive.ObjectID(id).Hex()
}

// Method to check if CustomObjectID is valid
func (id MongoID) IsValid() bool {
	return primitive.ObjectID(id).IsZero()
}
