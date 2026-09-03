package repository

import (
	"github.com/shappy0/mantritva/pkg/config"
	"gorm.io/gorm"
)

type Repos struct {
	DB                 *gorm.DB
	CassandraDB        *config.Cassandra
	MongoDB            *config.MongoDB
	Alliance           AllianceRepository
	Country            CountryRepository
	District           DistrictRepository
	Loksabha           LoksabhaRepository
	Member             MemberRepository
	Party              PartyRepository
	State              StateRepository
	Ruler              RulerRepository
	VidhanSabha        VidhansabhaRepository
	Leadership         LeadershipRepository
	NewsPressPost      NewsPressPostRepository
	User               UserRepository
	UserRole           UserRoleRepository
	Feed               FeedRepository
	Entity             EntityRepository
	Grievance          GrievanceRepository
	GrievanceCosign    GrievanceCosignRepository
	Report             ReportRepository
	Event              EventRepository
	DevelopmentWelfare DevelopmentWelfareRepository
	Enquiry            EnquiryRepository
	Reaction           ReactionRepository
	AnalyticStat       AnalyticStatRepository
	LocalBody          LocalBodyRepository
}

func NewRepos(db *gorm.DB, mongoDB *config.MongoDB) *Repos {
	return &Repos{
		DB: db,
		// CassandraDB:   cassandra,
		MongoDB:            mongoDB,
		Alliance:           NewAllianceRepository(db),
		Country:            NewCountryRepository(db),
		District:           NewDistrictRepository(db),
		Loksabha:           NewLoksabhaRepository(db),
		Member:             NewMemberRepository(db),
		Party:              NewPartyRepository(db),
		State:              NewStateRepository(db),
		Ruler:              NewRulersRepository(db),
		VidhanSabha:        NewVidhansabhaRepository(db),
		Leadership:         NewLeadershipRepository(db),
		NewsPressPost:      NewNewsPressPostRepository(db, mongoDB),
		User:               NewUserRepository(db),
		Feed:               NewFeedRepository(mongoDB),
		Entity:             NewEntityRepository(db),
		Grievance:          NewGrievanceRepository(db, mongoDB),
		GrievanceCosign:    NewGrievanceCosignRepository(mongoDB),
		Report:             NewReportRepository(db, mongoDB),
		Event:              NewEventRepository(mongoDB),
		DevelopmentWelfare: NewDevelopmentWelfareRepository(db, mongoDB),
		Enquiry:            NewEnquiryRepository(db),
		Reaction:           NewReactionRepository(mongoDB),
		AnalyticStat:       NewAnalyticStat(mongoDB),
		LocalBody:          NewLocalBodyRepository(db),
		UserRole:           NewUserRoleRepository(db),
	}
}
