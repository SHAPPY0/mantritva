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
type Vote string
type Position string
type GrievanceStatus string
type KRA string
type DevelopmentStatus string
type GrievanceSeverity string
type OTPMode string
type OTPPurpose string
type LocalBodyType string
type ULBType string

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
	return fmt.Errorf("Cannot convert %v to DateOnly", value)
}

func (d *DateOnly) IsNull() bool {
	return d == nil || d.Time == nil || d.Time.IsZero()
}

func (d DateOnly) String() string {
	if d.Time == nil || d.Time.IsZero() {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

// Method to convert CustomObjectID to string
func (id MongoID) String() string {
	return primitive.ObjectID(id).Hex()
}

// Method to check if CustomObjectID is valid
func (id MongoID) IsValid() bool {
	return primitive.ObjectID(id).IsZero()
}

func (p Position) Label() string {
	if val, ok := PositionDisplayMap[p]; ok {
		return val
	}
	return string(p)
}

func (kra KRA) Label() string {
	if val, ok := KraDisplayMap[kra]; ok {
		return val
	}
	return string(kra)
}

func (ulb ULBType) Label() string {
	if val, ok := ULBTypeDisplayMap[ulb]; ok {
		return val
	}
	return string(ulb)
}
