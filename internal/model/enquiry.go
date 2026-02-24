package model

type Enquiry struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Mobile      string   `json:"mobile"`
	EnquiryType string   `json:"enquiry_type"`
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
	FileUrls    []string `json:"file_urls" gorm:"type:jsonb;serializer:json"`
	AuditFields
}
