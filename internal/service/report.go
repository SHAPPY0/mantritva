package service

import (
	"math"
	"runtime"
	"sync"
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
		AuditFieldsBSON: model.InsertBSONAuditFields(""),
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
	pagination := dto.Pagination{Offset: 0, Limit: 500}
	newsPresses, err := s.NewsPressRepo.GetPostsByMemberID(pagination, memberID)
	if err != nil {
		return nil, err
	}
	report.Total = uint(len(*newsPresses))
	if report.Total == 0 {
		return &report, nil
	}
	maxWorkers := runtime.NumCPU()
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalNewsPressScore float64
	for _, eachNewsPress := range *newsPresses {
		wg.Add(1)
		sem <- struct{}{}
		p := eachNewsPress
		go func(post model.NewsPressPost) {
			defer wg.Done()
			defer func() { <-sem }()

			reactions, err := s.ReportRepo.GetReactionsByPostID(post.ID)
			if err != nil {
				log.Errorf("Error fetching reaction of newsPress %d", post.ID)
				return
			}
			if len(reactions) == 0 {
				return
			}
			eachScoreSum := 0
			for _, reaction := range reactions {
				eachScoreSum += reaction.NormalizedScore
			}
			//get average of each post's reaction
			avg := float64(eachScoreSum) / float64(len(reactions))
			mu.Lock()
			totalNewsPressScore += avg
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	// average across all posts (including posts with zero reactions as 0)
	rawNewsPressScore := totalNewsPressScore / float64(report.Total)
	if rawNewsPressScore < 0 {
		rawNewsPressScore = 0
	}
	if rawNewsPressScore > 100 {
		rawNewsPressScore = 100
	}
	//cover to scale 0 - weight
	scaled := (rawNewsPressScore / 100.0) * float64(weight)
	report.Score = uint(math.Round(scaled))
	return &report, nil
}

func (s *reportService) GetReport(memberID int) (*model.MemberReport, error) {
	currentYear := time.Now().Format("2006")
	report, err := s.ReportRepo.GetMemberScore(uint(memberID), currentYear)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, nil
	}

	return report, nil
}

func (svc *reportService) RunFeedScoreComputeJob(queue <-chan dto.ReactionReq) {
	for reactionReq := range queue {
		if reactionReq.FeedType == utils.NEWSPRESS {
			if err := svc.ComputeReactionScore(reactionReq); err != nil {
				log.Errorf(err.Error())
			}
		}
	}
}

func (svc *reportService) ComputeReactionScore(reaction dto.ReactionReq) error {
	var reactions []model.PostReaction
	var questionnaire *[]model.Questionnaire
	var err error
	if reactions, err = svc.ReportRepo.GetReactionsByPostID(reaction.PostID); err != nil {
		return err
	}
	if questionnaire, err = svc.ReportRepo.GetQuestionnairesByCategory(string(utils.NEWSPRESS)); err != nil {
		return err
	}
	totalVotes := len(reactions)
	totalWeight := 0
	for _, reaction := range reactions {
		weight := getWeight(questionnaire, reaction)
		if reaction.Reaction == string(utils.UPVOTE) {
			totalWeight += (1 + weight)
		} else if reaction.Reaction == string(utils.DOWNVOTE) {
			totalWeight -= (-1 - (weight))
			if totalWeight < 0 {
				totalWeight = 0
			}
		}
	}
	log.Infof("Total Votes %d, score;%d", totalVotes, totalWeight)
	// MIN_POSSIBLE := -10
	// MAX_POSSIBLE := 10
	// normalizedScore := ((totalWeight - MIN_POSSIBLE) / (MAX_POSSIBLE - MIN_POSSIBLE)) * 100
	return nil
}
