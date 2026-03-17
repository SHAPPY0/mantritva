package compute

import (
	"math"
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ─────────────────────────────────────────────────────────────────────────────
// GRIEVANCES
// ─────────────────────────────────────────────────────────────────────────────

// ComputeGrievanceScore calculates the Grievances KRA score.
//
// Grievances are different from posts — they have STATUS and SEVERITY that
// directly reflect the member's governance behaviour. The formula:
//
//	grievanceScore(g) = (statusWeight × 0.70 × 100) + (feedbackScore × 0.30 × 100)
//	                    weighted by severityWeight
//
//	finalScore = weightedSum / totalSeverity
//
// Penalties applied BEFORE Bayesian:
//   - Any HIGH/CRITICAL grievance with NO_UPDATE or OPEN → 30% score penalty
//   - 3+ severe unresolved grievances AND score > 40 → cap at 40
//   - Avg feedback < 0.25 on resolved grievances → 25% penalty
//   - Avg feedback < 0.10 on resolved grievances AND score > 50 → cap at 50
//
// Bayesian n = totalSeverity weight (not count), because 5 CRITICAL grievances
// represent more governance evidence than 5 LOW grievances.
func (svc Score) ComputeGrievanceScore(weight uint, memberID int, computeYear int, useDefault bool) (*model.GrievanceScore, error) {
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
		GlobalAverage: globalAverage,
		GrievanceIds:  make([]string, 0),
		IsProvisional: useDefault,
	}

	pagination := dto.Pagination{Offset: 0, Limit: maxPostsPerFetch}
	grievances, err := svc.GrievanceRepo.GetByMemberIDAndYear(pagination, memberID, computeYear)
	if err != nil {
		return nil, err
	}
	total := len(*grievances)
	kraReport.TotalCount = uint(total)

	if total == 0 {
		kraReport.Score = applyKRAWeight(globalAverage, weight)
		kraReport.Confidence = 0
		kraReport.IsProvisional = true
		kraReport.LastComputedAt = time.Now()
		return &kraReport, nil
	}

	var (
		totalSeverity         float64
		weightedSum           float64
		severeUnresolvedCount int
		resolvedCount         int
		totalReactionScore    float64
		totalLockedReactions  int
	)

	for _, grievance := range *grievances {
		// Skip grievances still under review or flagged invalid
		if grievance.FinalStatus == utils.UNDERREVIEW || grievance.FinalStatus == utils.INVALID {
			continue
		}

		severityW, ok := severityWeightage[grievance.Severity]
		if !ok {
			log.Warnf("ComputeGrievanceScore: unknown severity %q for grievance %s — skipping", grievance.Severity, grievance.GrievanceID)
			continue
		}
		statusW, ok := grievanceStatusWeightage[grievance.FinalStatus]
		if !ok {
			log.Warnf("ComputeGrievanceScore: unknown status %q for grievance %s — skipping", grievance.FinalStatus, grievance.GrievanceID)
			continue
		}

		// Severe neglect: HIGH/CRITICAL with no movement → stronger penalty
		if (grievance.Severity == utils.GRIEVANCE_HIGH || grievance.Severity == utils.GRIEVANCE_CRITICAL) &&
			(grievance.FinalStatus == utils.NOUPDATE || grievance.FinalStatus == utils.OPEN) {
			statusW = 0.1
			severeUnresolvedCount++
		}

		// 70% of grievance score comes from status weight
		grievanceScore := statusW * 0.70 * 100

		// 30% comes from public feedback — only applicable on resolved/closed
		if grievance.FinalStatus == utils.RESOLVED || grievance.FinalStatus == utils.CLOSED {
			feedbackScore := svc.computeGrievanceFeedback(grievance.ID)
			grievanceScore += feedbackScore * 0.30 * 100

			resolvedCount++
			totalReactionScore += feedbackScore
		}

		// Track locked reaction count for Bayesian n (even for non-resolved)
		if locked, err := svc.ReactionRepo.GetLockedReactions(grievance.ID); err == nil {
			totalLockedReactions += len(locked)
		}

		weightedSum += grievanceScore * severityW
		totalSeverity += severityW
	}

	var finalScore float64
	if totalSeverity == 0 {
		// All grievances were UNDERREVIEW/INVALID — treat as no data
		finalScore = globalAverage
	} else {
		finalScore = weightedSum / totalSeverity
	}

	finalScore = math.Max(0, math.Min(100, finalScore))

	// Bayesian smoothing — n = totalSeverity (severity-weighted evidence mass)
	n := totalSeverity
	bayesScore := (n*finalScore + bayesianPriorStrength*globalAverage) / (n + bayesianPriorStrength)

	// ── Penalty A: severe neglect ──────────────────────────────────────────
	if severeUnresolvedCount > 0 {
		bayesScore *= 0.7
	}
	if severeUnresolvedCount >= 3 && bayesScore > 40 {
		bayesScore = 40
	}

	// ── Penalty B: horrible public feedback on resolved grievances ─────────
	if resolvedCount > 0 {
		avgFeedback := totalReactionScore / float64(resolvedCount)
		if avgFeedback < 0.25 {
			bayesScore *= 0.75
		}
		if avgFeedback < 0.10 && bayesScore > 50 {
			bayesScore = 50
		}
	}

	bayesScore = math.Max(0, math.Min(100, bayesScore)) // re-clamp after penalties

	kraReport.Score = applyKRAWeight(bayesScore, weight)
	kraReport.Confidence = confidenceScore(n)
	if totalLockedReactions < minEvidenceReactions {
		kraReport.IsProvisional = true
	}
	kraReport.LastComputedAt = time.Now()
	return &kraReport, nil
}

// computeGrievanceFeedback returns a 0-1 normalised feedback score for a
// resolved/closed grievance. Returns 0 on error or no reactions — callers
// treat missing feedback as neutral-negative (no feedback = not validated).
func (svc Score) computeGrievanceFeedback(postID primitive.ObjectID) float64 {
	reactions, err := svc.ReactionRepo.GetLockedReactions(postID)
	if err != nil {
		log.Errorf("computeGrievanceFeedback: error fetching reactions for grievance %v: %v", postID, err)
		return 0
	}
	if len(reactions) == 0 {
		log.Warnf("computeGrievanceFeedback: no locked reactions for grievance %v", postID)
		return 0
	}

	var scoreSum float64
	for _, r := range reactions {
		// Clamp here only because grievance feedback score is used as a ratio (0-1).
		// A -5 NormalizedScore on a grievance reaction would push feedback below 0,
		// which would invert the 30% feedback bonus into a penalty — unintended.
		// Grievance reactions represent satisfaction, not disapproval weighting.
		normalised := math.Max(0, math.Min(100, float64(r.NormalizedScore)))
		scoreSum += normalised
	}
	return (scoreSum / float64(len(reactions))) / 100.0
}
