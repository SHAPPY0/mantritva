package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/shappy0/mantritva/pkg/utils"
)

type Member struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	Name             string           `json:"name"`
	FullName         string           `json:"full_name"`
	DOB              utils.DateOnly   `json:"dob" gorm:"type:date"`
	Gender           string           `json:"gender"`
	Photo            string           `json:"photo"`
	Qualifications   string           `json:"qualifications"`
	MaritalStatus    string           `json:"marital_status"`
	FinancialAssets  string           `json:"financial_assets"`
	SocialMediaLinks SocialMediaLinks `json:"social_media_links" gorm:"type:json"`
	AuditFields
}

type SocialMediaLinks map[string]string

func (s *SocialMediaLinks) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

func (s SocialMediaLinks) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type CriminalRecord struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	MemberID       uint           `json:"member_id"`
	Offense        string         `json:"offense"`
	Status         string         `json:"status"`
	ConvictionDate utils.DateOnly `json:"conviction_date" gorm:"type:date"`
	Sentence       string         `json:"sentence"`
	AuditFields
}

type Career struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	MemberID      uint           `json:"member_id"`
	PartyID       uint           `json:"party_id"`
	Position      string         `json:"position"`
	TermStartDate utils.DateOnly `json:"term_start_date" gorm:"type:date"`
	TermEndDate   utils.DateOnly `json:"term_end_date,omitempty" gorm:"type:date"`
	AuditFields
}
