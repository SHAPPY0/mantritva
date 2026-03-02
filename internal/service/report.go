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
	GenerateReport(int) (*model.Score, error)
	GetReport(int) (dto.MemberReport, error)
}

type reportService struct {
	ReportRepo      repository.ReportRepository
	GrievanceRepo   repository.GrievanceRepository
	DevelopmentRepo repository.DevelopmentWelfareRepository
	NewsPressRepo   repository.NewsPressPostRepository
	PartyRepo       repository.PartyRepository
	ReactionRepo    repository.ReactionRepository
}

func NewReportService(repos *repository.Repos) ReportService {
	return &reportService{
		ReportRepo:      repos.Report,
		GrievanceRepo:   repos.Grievance,
		DevelopmentRepo: repos.DevelopmentWelfare,
		NewsPressRepo:   repos.NewsPressPost,
		PartyRepo:       repos.Party,
		ReactionRepo:    repos.Reaction,
	}
}

func (s *reportService) GenerateReport(memberID int) (*model.Score, error) {
	scoreData := model.Score{
		ID:              primitive.NewObjectID(),
		MemberID:        uint(memberID),
		MaxScore:        100,
		Year:            time.Now().Format("2006"),
		Version:         1,
		AuditFieldsBSON: model.InsertBSONAuditFields(""),
	}
	var ksaWeights dto.KRAWeights
	kras, err := s.ReportRepo.GetReportKRA()
	if err != nil {
		return nil, err
	}

	for _, kra := range kras {
		if kra.KRAName == string(utils.KRA_DEVELOPMENTS) {
			ksaWeights.Developments = kra.Weight
		} else if kra.KRAName == string(utils.KRA_GRIEVANCES) {
			ksaWeights.Grievances = kra.Weight
		} else if kra.KRAName == string(utils.KRA_CALRECORDS) {
			ksaWeights.CALRecords = kra.Weight
		} else if kra.KRAName == string(utils.KRA_NEWSPRESSES) {
			ksaWeights.NewsPresses = kra.Weight
		}
	}

	var useDefault bool
	startDate, err := s.PartyRepo.GetPoliticalStartDate(uint(memberID))
	if err != nil {
		log.Errorf("Error while fetching political start date for memberID %d: %v", memberID, err)
	} else {
		checkPoliticalStartDate, _ := time.Parse("2006-01-02", startDate.Format("2006-01-02"))
		if time.Since(checkPoliticalStartDate).Hours() > 24*30*6 { // If member has started political career within last 6 months then we can assume that they don't have enough data to calculate their score and we can use default score based on global averaget
			useDefault = true
		}
	}

	dReport, err := s.CalculateDevelopmentScore(ksaWeights.Developments, memberID, useDefault)
	if err != nil {
		log.Errorf("Error while fetching development&Welfare reports for memberID %d: %v", memberID, err)
	} else {
		score := dReport.Score
		//If score is zero
		if score == 0 {
			// And Member don't have historic data then set default score using Bayesian average, this is to ensure that new members are not at disadvantage and they get a fair score based on global average until we have enough data to calculate their actual score
			if useDefault {
				defaultScore, confidence := setDefaultScore(ksaWeights.Developments, dReport.GlobalAverage)
				dReport.Score = uint(defaultScore)
				dReport.Confidence = uint(confidence)
			}
		}
		scoreData.TotalScore += score
		scoreData.Confidence += dReport.Confidence
		scoreData.Developments = *dReport
	}

	gReport, err := s.CalculateGrievanceScore(ksaWeights.Grievances, memberID)
	if err != nil {
		log.Errorf("Error while fetching grievance reports for memberID %d: %v", memberID, err)
	} else {
		score := gReport.Score
		//If score is zero
		if score == 0 {
			// And Member don't have historic data then set default score using Bayesian average, this is to ensure that new members are not at disadvantage and they get a fair score based on global average until we have enough data to calculate their actual score
			if useDefault {
				defaultScore, confidence := setDefaultScore(ksaWeights.Grievances, gReport.GlobalAverage)
				gReport.Score = uint(defaultScore)
				gReport.Confidence = uint(confidence)
			}
		}
		scoreData.TotalScore += score
		scoreData.Confidence += gReport.Confidence
		scoreData.Grievances = *gReport
	}

	nsReport, err := s.CalculateNewsPressesScore(ksaWeights.NewsPresses, memberID)
	if err != nil {
		log.Errorf("Error while fetching newspresses reports for memberID %d: %v", memberID, err)
	} else {
		score := nsReport.Score
		//If score is zero
		if score == 0 {
			// And Member don't have historic data then set default score using Bayesian average, this is to ensure that new members are not at disadvantage and they get a fair score based on global average until we have enough data to calculate their actual score
			if useDefault {
				defaultScore, confidence := setDefaultScore(ksaWeights.NewsPresses, nsReport.GlobalAverage)
				nsReport.Score = uint(defaultScore)
				nsReport.Confidence = uint(confidence)
			}
		}
		scoreData.TotalScore += score
		scoreData.Confidence += nsReport.Confidence
		scoreData.NewsPresses = *nsReport
	}

	if err := s.ReportRepo.Create(&scoreData); err != nil {
		return nil, err
	}
	return &scoreData, nil
}

func (s *reportService) CalculateDevelopmentScore(weight uint, memberID int, useDefault bool) (*model.DevelopmentScore, error) {
	kraReport := model.DevelopmentScore{
		Score:          0,
		Weight:         weight,
		GlobalAverage:  50, //Global average can be calculated from historical data, for now we can assume it to be 50
		DevelopmentIds: make([]string, 0),
	}

	const bayesianWeight = 30.0

	developments, err := s.DevelopmentRepo.GetByMemberID(memberID)
	if err != nil {
		return nil, err
	}
	kraReport.TotalCount = uint(len(developments))
	if kraReport.TotalCount == 0 {
		return &kraReport, nil
	}
	maxWorkers := runtime.NumCPU()
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalWeightedScore float64
	var totalPostsCount float64

	for _, development := range developments {
		wg.Add(1)
		sem <- struct{}{}
		d := development
		go func(_development model.DevelopmentWelfare) {
			defer wg.Done()
			defer func() { <-sem }()
			reactions, err := s.ReactionRepo.GetByPostID(d.ID)
			if err != nil {
				log.Errorf("Error fetching reaction of development %d", d.ID)
				return
			}
			reactionCount := float64(len(reactions))
			if reactionCount == 0 {
				return
			}
			var scoreSum float64
			for _, reaction := range reactions {
				scoreSum += float64(reaction.NormalizedScore)
			}
			//get average of each post's reaction
			postAverage := scoreSum / reactionCount
			// Bayesian smoothing
			adjustedPostScore := (reactionCount*postAverage + bayesianWeight*float64(kraReport.GlobalAverage)) / (reactionCount + bayesianWeight)

			mu.Lock()
			totalWeightedScore += adjustedPostScore
			totalPostsCount += 1
			mu.Unlock()
		}(d)
	}
	wg.Wait()
	// average across all posts (including posts with zero reactions as 0)
	rawDevelopmentScore := totalWeightedScore / totalPostsCount
	if rawDevelopmentScore < 0 {
		rawDevelopmentScore = 0
	}
	if rawDevelopmentScore > 100 {
		rawDevelopmentScore = 100
	}
	//cover to scale 0 - weight
	scaled := (rawDevelopmentScore / 100.0) * float64(weight)
	kraReport.Score = uint(math.Round(scaled))
	kraReport.Confidence = uint(math.Min(100, (totalPostsCount/(totalPostsCount+bayesianWeight))*100))
	kraReport.LastActivityAt = time.Now()
	return &kraReport, nil
}

func (s *reportService) CalculateGrievanceScore(weight uint, memberID int) (*model.GrievanceScore, error) {
	kraReport := model.GrievanceScore{
		Score:         0,
		Weight:        weight,
		GlobalAverage: 50.0, //Global average can be calculated from historical data, for now we can assume it to be 50
		GrievanceIds:  make([]string, 0),
	}

	pagination := dto.Pagination{Offset: 0, Limit: 500}
	grievances, err := s.GrievanceRepo.GetByMemberID(pagination, memberID, "")
	if err != nil {
		return nil, err
	}
	total := len(*grievances)
	kraReport.TotalCount = uint(total)

	if total == 0 {
		return &kraReport, nil
	}
	var closed float64
	for _, grievance := range *grievances {
		if grievance.FinalStatus == utils.COMPLETED || grievance.FinalStatus == utils.FIXED {
			closed++
		} else {
			kraReport.GrievanceIds = append(kraReport.GrievanceIds, grievance.GrievanceID)
		}
	}
	resolutionRate := closed / float64(total)
	// Bayesian smoothing
	bayesianWeight := 30.0
	adjustedScore := ((float64(total) * resolutionRate) + (bayesianWeight * float64(kraReport.GlobalAverage))) / (float64(total) + bayesianWeight)

	kraReport.Score = uint(float64(weight) * adjustedScore)
	kraReport.Confidence = uint(math.Min(100, (float64(total)/(float64(total)+bayesianWeight))*100))
	kraReport.LastActivityAt = time.Now()
	return &kraReport, nil
}

func (s *reportService) CalculateNewsPressesScore(weight uint, memberID int) (*model.NewsPressScore, error) {
	kraReport := model.NewsPressScore{
		Score:         0,
		Weight:        weight,
		GlobalAverage: 50.0, //Global average can be calculated from historical data, for now we can assume it to be 50
		NewsPressIds:  make([]string, 0),
	}
	pagination := dto.Pagination{Offset: 0, Limit: 500}
	newsPresses, err := s.NewsPressRepo.GetPostsByMemberID(pagination, memberID)
	if err != nil {
		return nil, err
	}
	kraReport.TotalCount = uint(len(*newsPresses))
	if kraReport.TotalCount == 0 {
		return &kraReport, nil
	}
	maxWorkers := runtime.NumCPU()
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	const bayesianWeight = 30.0
	var totalWeightedScore float64
	var totalPostsCount float64

	for _, post := range *newsPresses {
		wg.Add(1)
		sem <- struct{}{}
		p := post
		go func(_post model.NewsPressPost) {
			defer wg.Done()
			defer func() { <-sem }()

			reactions, err := s.ReactionRepo.GetByPostID(_post.ID)
			if err != nil {
				log.Errorf("Error fetching reaction of newsPress %d", _post.ID)
				return
			}
			reactionCount := float64(len(reactions))
			if reactionCount == 0 {
				return
			}
			var scoreSum float64
			for _, reaction := range reactions {
				scoreSum += float64(reaction.NormalizedScore)
			}
			//get average of each post's reaction
			postAverage := scoreSum / reactionCount
			// Bayesian smoothing
			adjustedPostScore := (reactionCount*postAverage + bayesianWeight*float64(kraReport.GlobalAverage)) / (reactionCount + bayesianWeight)

			mu.Lock()
			totalWeightedScore += adjustedPostScore
			totalPostsCount += 1
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	// average across all posts (including posts with zero reactions as 0)
	rawNewsPressScore := totalWeightedScore / totalPostsCount
	if rawNewsPressScore < 0 {
		rawNewsPressScore = 0
	}
	if rawNewsPressScore > 100 {
		rawNewsPressScore = 100
	}
	//cover to scale 0 - weight
	scaled := (rawNewsPressScore / 100.0) * float64(weight)
	kraReport.Score = uint(math.Round(scaled))
	kraReport.Confidence = uint(math.Min(100, (totalPostsCount/(totalPostsCount+bayesianWeight))*100))
	kraReport.LastActivityAt = time.Now()
	return &kraReport, nil
}

func setDefaultScore(kraWeight uint, globalAvg uint) (float64, float64) {
	//using Bayesian formula
	score := (float64(globalAvg) / 100) * float64(kraWeight)
	return score, 0
}

func (s *reportService) GetReport(memberID int) (dto.MemberReport, error) {
	currentYear := time.Now().Format("2006")
	memberReport := dto.MemberReport{MemberID: uint(memberID), Year: currentYear}
	report, err := s.ReportRepo.GetMemberScore(uint(memberID), currentYear)
	if err != nil {
		return memberReport, err
	}
	if report == nil {
		return memberReport, nil
	}
	memberReport.TotalScore = report.TotalScore
	memberReport.CreatedAt = report.CreatedAt
	memberReport.KRAS = dto.KRAReport{
		Developments: dto.DevelopmentReport{
			Name:         string(utils.KRA_DEVELOPMENTS),
			Score:        report.Developments.Score,
			Weight:       report.Developments.Weight,
			Total:        report.Developments.Total,
			Developments: report.Developments.Developments,
		},
		Grievances: dto.GrievanceReport{
			Name:       string(utils.KRA_GRIEVANCES),
			Score:      report.Grievances.Score,
			Weight:     report.Grievances.Weight,
			Total:      report.Grievances.Total,
			Grievances: report.Grievances.Grievances,
		},
		NewsPresses: dto.NewsPressesReport{
			Name:        string(utils.KRA_NEWSPRESSES),
			Score:       report.NewsPresses.Score,
			Weight:      report.NewsPresses.Weight,
			Total:       report.NewsPresses.Total,
			NewsPresses: report.NewsPresses.NewsPresses,
		},
		CALRecords: dto.CALReport{
			Name:       string(utils.KRA_CALRECORDS),
			Score:      report.CALRecords.Score,
			Weight:     report.CALRecords.Weight,
			Total:      report.CALRecords.Total,
			CALRecords: report.CALRecords.CALRecords,
		},
	}
	return memberReport, nil
}
