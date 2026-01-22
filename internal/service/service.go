package service

import "github.com/shappy0/mantritva/internal/repository"

type Services struct {
	NewsPressService NewsPressService
	FeedService      FeedService
	ConsoleService   ConsoleService
	GrievanceService GrievanceSerivce
	ReportService    ReportService
	EventService     EventService
}

func NewServices(repos *repository.Repos) *Services {
	return &Services{
		NewsPressService: NewNewsPressService(repos),
		FeedService:      NewFeedService(repos),
		ConsoleService:   NewConsoleService(repos),
		GrievanceService: NewGrievanceService(repos),
		ReportService:    NewReportService(repos),
		EventService:     NewEventService(repos),
	}
}
