package dto

type ConsoleDashboardRes struct {
	Users               int64 `json:"users"`
	NewsPress           int64 `json:"news_press"`
	DevelopmentWelfares int64 `json:"development_welfares"`
	Grievances          int64 `json:"grievances"`
	Total               int64 `json:"total"`
}
