package service

import (
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/cron"
)

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
	CronJob                   *cron.CronJob
	AnalyticStatsService      AnalyticStatsService
	LocalBodyService          *LocalBodyService
	SuggestionService         SuggestionService
}

func NewServices(repos *repository.Repos) *Services {
	svcs := &Services{
		UserService:               NewUserService(repos),
		NewsPressService:          NewNewsPressService(repos),
		FeedService:               NewFeedService(repos),
		ConsoleService:            NewConsoleService(repos),
		MemberService:             NewMemberService(repos),
		GrievanceService:          NewGrievanceService(repos),
		EventService:              NewEventService(repos),
		DevelopmentWelfareService: NewDevelopmentWelfareService(repos),
		CommonService:             NewCommonService(repos),
		AnalyticStatsService:      NewAnalyticStatService(repos),
		LocalBodyService:          NewLocalBodyService(repos),
		SuggestionService:         NewSuggestionService(repos),
	}
	svcs.ReportService = NewReportService(repos, svcs)
	return svcs
}
