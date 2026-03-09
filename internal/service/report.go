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
	GetMemberReports(int, int, int) ([]dto.MemberReport, error)
	// GetMemberReportCard(string) (dto.MemberReportCard, error)
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
	var totalScore float64
	var totalWeight float64
	var totalConfidenceWeighted float64
	currentYear := time.Now().Year()
	scoreData := model.Score{
		ID:              primitive.NewObjectID(),
		MemberID:        uint(memberID),
		MaxScore:        100,
		Year:            currentYear,
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
				dReport.Score = defaultScore
				dReport.Confidence = confidence
			}
		}
		totalScore += dReport.Score
		totalWeight += float64(ksaWeights.Developments)
		totalConfidenceWeighted += dReport.Confidence * float64(ksaWeights.Developments)
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
				gReport.Score = defaultScore
				gReport.Confidence = confidence
			}
		}
		totalScore += gReport.Score
		totalWeight += float64(ksaWeights.Grievances)
		totalConfidenceWeighted += gReport.Confidence * float64(ksaWeights.Grievances)
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
				nsReport.Score = defaultScore
				nsReport.Confidence = confidence
			}
		}
		totalScore += nsReport.Score
		totalWeight += float64(ksaWeights.NewsPresses)
		totalConfidenceWeighted += nsReport.Confidence * float64(ksaWeights.NewsPresses)
		scoreData.NewsPresses = *nsReport
	}

	if totalWeight > 0 {
		scoreData.TotalScore = (totalScore / totalWeight) * 100 // convert to percentage
		scoreData.Confidence = totalConfidenceWeighted / totalWeight
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
				normalizedScore := math.Max(0, math.Min(100, float64(reaction.NormalizedScore)))
				scoreSum += normalizedScore
			}
			//get average of each post's reaction
			postAverage := scoreSum / reactionCount
			// Bayesian smoothing
			// bayesianScore := (reactionCount*postAverage + bayesianWeight*float64(kraReport.GlobalAverage)) / (reactionCount + bayesianWeight)

			mu.Lock()
			totalWeightedScore += postAverage
			totalPostsCount += 1
			mu.Unlock()
		}(d)
	}
	wg.Wait()
	// average across all posts (including posts with zero reactions as 0)
	var rawDevelopmentScore float64
	if totalPostsCount > 0 {
		rawDevelopmentScore = totalWeightedScore / totalPostsCount
	}

	//Apply Bayesian smoothing Once at the aggregate level
	//This balances between the calculated score and the gloabl average.
	if totalPostsCount+bayesianWeight > 0 {
		bayesianScore := ((totalPostsCount * rawDevelopmentScore) + (bayesianWeight * float64(kraReport.GlobalAverage))) / (totalPostsCount + bayesianWeight)
		rawDevelopmentScore = bayesianScore
	}

	if rawDevelopmentScore < 0 {
		rawDevelopmentScore = 0
	}
	if rawDevelopmentScore > 100 {
		rawDevelopmentScore = 100
	}
	//cover to scale 0 - weight
	scaled := (rawDevelopmentScore / 100.0) * float64(weight)
	kraReport.Score = math.Round(scaled)
	if (totalPostsCount + bayesianWeight) == 0 {
		kraReport.Confidence = 5.0
	} else {
		kraReport.Confidence = math.Min(100, (totalPostsCount/(totalPostsCount+bayesianWeight))*100)
	}
	kraReport.LastActivityAt = time.Now()
	return &kraReport, nil
}

func (s *reportService) CalculateGrievanceScore(weight uint, memberID int) (*model.GrievanceScore, error) {
	severityWeightage := map[utils.GrievanceSeverity]float64{
		utils.GRIEVANCE_LOW:      1.0,
		utils.GRIEVANCE_MEDIUM:   1.5,
		utils.GRIEVANCE_HIGH:     2.0,
		utils.GRIEVANCE_CRITICAL: 2.5,
	}
	grievanceStatusWeightage := map[utils.GrievanceStatus]float64{
		utils.RESOLVED:   1.0,
		utils.CLOSED:     0.8,
		utils.INPROGRESS: 0.5,
		utils.OPEN:       0.3,
		utils.STAYED:     0.5,
		utils.NOUPDATE:   0.2,
		utils.REJECTED:   0.25,
	}
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
	var totalSeverity float64
	var weightedSum float64
	var severeUnresolvedCount int
	var resolvedCount int
	var totalReactionScore float64

	for _, grievance := range *grievances {
		if grievance.FinalStatus == utils.UNDERREVIEW || grievance.FinalStatus == utils.INVALID {
			continue
		}
		//severity weight
		severityW, ok := severityWeightage[grievance.Severity]
		if !ok {
			continue
		}
		//status weight
		statusW, ok := grievanceStatusWeightage[grievance.FinalStatus]
		if !ok {
			continue
		}
		// If severity HIGH or CRITICAL AND status NO_UPDATE → stronger penalty
		if (grievance.Severity == utils.GRIEVANCE_HIGH || grievance.Severity == utils.GRIEVANCE_CRITICAL) &&
			(grievance.FinalStatus == utils.NOUPDATE || grievance.FinalStatus == utils.OPEN) {
			statusW = 0.1
			severeUnresolvedCount++
		}

		//BaseScore becomes equal to Ws after normalization logic
		baseScore := statusW
		// 70% weighted grievance score of baseScore (0–100 scale)
		grievanceScore := (baseScore * 0.70) * 100

		//Feedaback  weight
		if grievance.FinalStatus == utils.RESOLVED || grievance.FinalStatus == utils.CLOSED {
			var reactionScore float64
			var scoreSum float64
			reactions, err := s.ReactionRepo.GetByPostID(grievance.ID)
			if err == nil {
				reactionCount := float64(len(reactions))
				if reactionCount == 0 {
					log.Warnf("No reaction found for grievance %d", grievance.ID)
					reactionScore = 0.0 //Default to 0 if no reactions
				} else {
					for _, reaction := range reactions {
						if reaction.NormalizedScore < 0 || reaction.NormalizedScore > 100 {
							log.Warnf("Invalid normalizeScore %d for reaction in grievance %d clamping to 0-100", reaction.NormalizedScore, grievance.ID)
						}
						normalizedScore := math.Max(0, math.Min(100, float64(reaction.NormalizedScore)))
						scoreSum += normalizedScore
					}
					//Calculate average feedback score (0-100) then normalize 0-1 for weight
					avgReactionScore := scoreSum / reactionCount
					reactionScore = avgReactionScore / 100.0 //Normalize to 0-1 scale
				}
			} else {
				reactionScore = 0.0 // Default to 0 on error
				log.Errorf("Error while fetching reaction of grievance %s", grievance.GrievanceID)
			}
			// Clamp reactionScore to 0-1 (though it should already be)
			reactionScore = math.Max(0, math.Min(1, reactionScore))
			resolvedCount++
			totalReactionScore += reactionScore
			// Rest 30% weighted grievance score of feedback (0–100 scale)
			// Use reactionScore (0-1) scaled to contribute 0-30 to grievanceScore
			grievanceScore += (reactionScore * 0.3) * 100
		}
		// Aggregate weighted by severity
		weightedSum += grievanceScore * severityW
		totalSeverity += severityW
	}

	var finalScore float64
	if totalSeverity == 0 {
		finalScore = float64(kraReport.GlobalAverage)
	} else {
		finalScore = weightedSum / totalSeverity
	}
	// RULE A: Severe Neglect Penalty
	if severeUnresolvedCount > 0 {
		finalScore *= 0.7
	}

	if severeUnresolvedCount >= 3 && finalScore > 40 {
		finalScore = 40
	}

	//Horrible Feedback Penalty
	if resolvedCount > 0 {
		avgReaction := totalReactionScore / float64(resolvedCount)
		//If avgReaction is less then 0.25 then add penalty. Reduce the grievance score by 25% when average user feedback for resolved grievances is very low
		if avgReaction < 0.25 {
			finalScore *= 0.75
		}
		if avgReaction < 0.10 && finalScore > 50 {
			finalScore = 50
		}
	}

	if finalScore < 0 {
		finalScore = 0
	}
	if finalScore > 100 {
		finalScore = 100
	}
	bayesianWeight := 30.0 // stabilization strength
	// Not grievance count? Because not all grievances are equal Severity already represents importance.5 CRITICAL grievances is stronger evidence than 5 LOW grievances
	// It means: “How much meaningful governance evidence do we have?” Not just “how many grievances handled".
	n := totalSeverity
	//prevent division by zero in Bayesian calculation
	if n+bayesianWeight == 0 {
		finalScore = float64(kraReport.GlobalAverage)
	} else {
		bayesianScore := ((n / (n + bayesianWeight)) * finalScore) + (float64(bayesianWeight) / float64(n+bayesianWeight) * float64(kraReport.GlobalAverage))
		finalScore = bayesianScore
	}

	scaled := (finalScore / 100.0) * float64(weight)
	kraReport.Score = math.Round(scaled)
	kraReport.Confidence = math.Min(100, (n/(n+bayesianWeight))*100)
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
	kraReport.Score = math.Round(scaled)
	kraReport.Confidence = math.Min(100, (totalPostsCount/(totalPostsCount+bayesianWeight))*100)
	kraReport.LastActivityAt = time.Now()
	return &kraReport, nil
}

func setDefaultScore(kraWeight uint, globalAvg float64) (float64, float64) {
	//using Bayesian formula
	score := (globalAvg / 100.0) * float64(kraWeight)
	return score, 5 //default to 5
}

func prepareReport(report *model.Score, memberReportResult *dto.MemberReport) {
	memberReportResult.TotalScore = utils.RoundToOneDecimal(report.TotalScore)
	memberReportResult.Confidence = utils.RoundToOneDecimal(report.Confidence)
	memberReportResult.CreatedAt = report.CreatedAt
	memberReportResult.KRAS = dto.KRAReport{
		Developments: dto.DevelopmentReport{
			Name:         string(utils.KRA_DEVELOPMENTS),
			Score:        utils.RoundToOneDecimal(report.Developments.Score),
			Confidence:   utils.RoundToOneDecimal(report.Developments.Confidence),
			Weight:       report.Developments.Weight,
			Total:        report.Developments.TotalCount,
			Developments: report.Developments.DevelopmentIds,
		},
		Grievances: dto.GrievanceReport{
			Name:       string(utils.KRA_GRIEVANCES),
			Score:      utils.RoundToOneDecimal(report.Grievances.Score),
			Confidence: utils.RoundToOneDecimal(report.Grievances.Confidence),
			Weight:     report.Grievances.Weight,
			Total:      report.Grievances.TotalCount,
			Grievances: report.Grievances.GrievanceIds,
		},
		NewsPresses: dto.NewsPressesReport{
			Name:        string(utils.KRA_NEWSPRESSES),
			Score:       utils.RoundToOneDecimal(report.NewsPresses.Score),
			Confidence:  utils.RoundToOneDecimal(report.NewsPresses.Confidence),
			Weight:      report.NewsPresses.Weight,
			Total:       report.NewsPresses.TotalCount,
			NewsPresses: report.NewsPresses.NewsPressIds,
		},
		CALRecords: dto.CALReport{
			Name:       string(utils.KRA_CALRECORDS),
			Score:      utils.RoundToOneDecimal(report.CALRecords.Score),
			Confidence: utils.RoundToOneDecimal(report.CALRecords.Confidence),
			Weight:     report.CALRecords.Weight,
			Total:      report.CALRecords.TotalCount,
			CALRecords: report.CALRecords.CALRecordIds,
		},
	}
}

func (s *reportService) GetReport(memberID int) (dto.MemberReport, error) {
	currentYear := time.Now().Year()
	memberReport := dto.MemberReport{MemberID: uint(memberID), Year: currentYear}
	report, err := s.ReportRepo.GetMemberScore(uint(memberID), currentYear)
	if err != nil {
		return memberReport, err
	}
	if report == nil {
		return memberReport, nil
	}
	prepareReport(report, &memberReport)
	return memberReport, nil
}

func (s *reportService) GetMemberReports(memberID int, from, to int) ([]dto.MemberReport, error) {
	memberReports := make([]dto.MemberReport, 0)
	reports, err := s.ReportRepo.GetMemberScoreReports(uint(memberID), from, to)
	if err != nil {
		return memberReports, nil
	}
	for _, report := range reports {
		memberReport := dto.MemberReport{MemberID: uint(memberID), Year: report.Year}
		prepareReport(report, &memberReport)
		memberReports = append(memberReports, memberReport)
	}
	return memberReports, nil
}
