package dto

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
