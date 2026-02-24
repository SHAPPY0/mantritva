package dto

type Enquiry struct {
	Name        string `form:"name"`
	Email       string `form:"email"`
	Mobile      string `form:"mobile"`
	EnquiryType string `form:"enquiry_type"`
	Subject     string `form:"subject"`
	Description string `form:"description"`
}
