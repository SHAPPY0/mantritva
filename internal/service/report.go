package service

import (
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService interface {
	GenerateReport(int) (*model.MemberReport, error)
	GetReport(int) (*model.MemberReport, error)
}

type reportService struct {
	ReportRepo    repository.ReportRepository
	GrievanceRepo repository.GrievanceRepository
	NewsPressRepo repository.NewsPressPostRepository
}

func NewReportService(repos *repository.Repos) ReportService {
	return &reportService{
		ReportRepo:    repos.Report,
		GrievanceRepo: repos.Grievance,
		NewsPressRepo: repos.NewsPressPost,
	}
}

func (s *reportService) GenerateReport(memberID int) (*model.MemberReport, error) {
	report := model.MemberReport{
		ID:              primitive.NewObjectID(),
		MemberID:        uint(memberID),
		MaxScore:        100,
		Year:            time.Now().Format("2006"),
		AuditFieldsBSON: model.InsertBSONAuditFields(),
	}
	var ksaWeights dto.KSAWeights
	ksas, err := s.ReportRepo.GetReportKSA()
	if err != nil {
		return nil, err
	}
	for _, ksa := range ksas {
		if ksa.KSAName == string(utils.KSA_DEVELOPMENTS) {
			ksaWeights.Developments = ksa.Weight
		} else if ksa.KSAName == string(utils.KSA_GRIEVANCES) {
			ksaWeights.Grievances = ksa.Weight
		} else if ksa.KSAName == string(utils.KSA_FIRS) {
			ksaWeights.FIRs = ksa.Weight
		} else if ksa.KSAName == string(utils.KSA_NEWSPRESSES) {
			ksaWeights.NewsPresses = ksa.Weight
		}
	}

	dReport, err := s.CalculateDevelopmentScore(ksaWeights.Developments, memberID)
	if err != nil {
		log.Errorf("Error while fetching grievance report for memberID %d: %v", memberID, err)
	} else {
		report.TotalScore += dReport.Score
		report.Developments = *dReport
	}

	gReport, err := s.CalculateGrievanceScore(ksaWeights.Grievances, memberID)
	if err != nil {
		log.Errorf("Error while fetching grievance report for memberID %d: %v", memberID, err)
	} else {
		report.TotalScore += gReport.Score
		report.Grievances = *gReport
	}

	nsReport, err := s.CalculateNewsPressesScore(ksaWeights.NewsPresses, memberID)
	if err != nil {
		log.Errorf("Error while fetching newspresses report for memberID %d: %v", memberID, err)
	} else {
		report.TotalScore += nsReport.Score
		report.NewsPresses = *nsReport
	}

	if err := s.ReportRepo.Create(&report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *reportService) CalculateDevelopmentScore(weight uint, memberID int) (*model.DevelopmentReport, error) {
	report := model.DevelopmentReport{
		Score:        weight,
		Weight:       weight,
		Developments: make([]string, 0),
	}
	return &report, nil
}

func (s *reportService) CalculateGrievanceScore(weight uint, memberID int) (*model.GrievanceReport, error) {
	report := model.GrievanceReport{
		Score:      weight,
		Weight:     weight,
		Grievances: make([]string, 0),
	}
	stillOpen := 0
	pagination := dto.Pagination{Offset: 0, Limit: 500}
	grievances, err := s.GrievanceRepo.GetGrievancesByMemberID(pagination, memberID, "")
	if err != nil {
		return nil, err
	}
	report.Total = uint(len(*grievances))
	if report.Total == 0 {
		return &report, nil
	}
	for _, grievance := range *grievances {
		if grievance.FinalStatus != utils.COMPLETED && grievance.FinalStatus != utils.FIXED {
			stillOpen += 1
			// report.Grievances = append(report.Grievances, dto.GrievanceShort{
			// 	GrievanceID: grievance.GrievanceID,
			// 	Title:       grievance.Title,
			// 	FinalStatus: grievance.FinalStatus,
			// })
			report.Grievances = append(report.Grievances, grievance.GrievanceID)
		}
	}
	report.Score = weight * (1 - uint(stillOpen)/uint(report.Total))
	return &report, nil
}

func (s *reportService) CalculateNewsPressesScore(weight uint, memberID int) (*model.NewsPressReport, error) {
	report := model.NewsPressReport{
		Score:       weight,
		Weight:      weight,
		NewsPresses: make([]string, 0),
	}
	notLikedPosts := 0
	pagination := dto.Pagination{Offset: 0, Limit: 500}
	newsPresses, err := s.NewsPressRepo.GetPostsByMemberID(pagination, memberID)
	log.Infof("News And Presses %+v", newsPresses)
	if err != nil {
		return nil, err
	}
	report.Total = uint(len(*newsPresses))
	if report.Total == 0 {
		return &report, nil
	}
	for _, ns := range *newsPresses {
		if ns.UpvoteCount < ns.DownvoteCount {
			notLikedPosts += 1
			// report.NewsPresses = append(report.NewsPresses, dto.NewsPressShort{
			// 	ID:            ns.ID,
			// 	Statement:     ns.Statement,
			// 	StatementDate: ns.StatementDate,
			// 	UpvoteCount:   ns.UpvoteCount,
			// 	DownvoteCount: ns.DownvoteCount,
			// })
			report.NewsPresses = append(report.NewsPresses, ns.ID.Hex())
		}
	}
	report.Score = weight * (1 - uint(notLikedPosts)/uint(report.Total))
	return &report, nil
}

func (s *reportService) GetReport(memberID int) (*model.MemberReport, error) {
	report, err := s.ReportRepo.GetMemberScore(uint(memberID), time.Now().Format("2006"))
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, nil
	}
	return report, nil
}
