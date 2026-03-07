package service

import "github.com/shappy0/mantritva/internal/repository"

type Services struct {
	UserService               UserService
	NewsPressService          NewsPressService
	FeedService               FeedService
	ConsoleService            ConsoleService
	GrievanceService          GrievanceSerivce
	MemberService             MemberService
	ReportService             ReportService
	EventService              EventService
	DevelopmentWelfareService DevelopmentWelfareService
	CommonService             CommonService
}

func NewServices(repos *repository.Repos) *Services {
	return &Services{
		UserService:               NewUserService(repos),
		NewsPressService:          NewNewsPressService(repos),
		FeedService:               NewFeedService(repos),
		ConsoleService:            NewConsoleService(repos),
		MemberService:             NewMemberService(repos),
		GrievanceService:          NewGrievanceService(repos),
		ReportService:             NewReportService(repos),
		EventService:              NewEventService(repos),
		DevelopmentWelfareService: NewDevelopmentWelfareService(repos),
		CommonService:             NewCommonService(repos),
	}
}
