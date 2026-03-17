package compute

import (
	"math"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
)

const (
	// bayesianPriorStrength controls how strongly the global average pulls
	// the score when evidence (locked reactions) is thin.
	// 30 means you need ~30 locked reactions before your real data dominates.
	bayesianPriorStrength = 30.0

	// globalAverage is the neutral prior score (0-100 scale).
	// Replace with a live DB aggregate once you have enough historical data.
	globalAverage = 50.0

	// minEvidenceReactions: below this count the KRA score is marked provisional.
	minEvidenceReactions = 0

	// maxPostsPerFetch caps the pagination window for post fetches.
	maxPostsPerFetch = 500
)

type Score struct {
	DevelopmentRepo repository.DevelopmentWelfareRepository
	NewsPressRepo   repository.NewsPressPostRepository
	GrievanceRepo   repository.GrievanceRepository
	ReactionRepo    repository.ReactionRepository
}

func NewScore(repos *repository.Repos) Score {
	return Score{
		DevelopmentRepo: repos.DevelopmentWelfare,
		NewsPressRepo:   repos.NewsPressPost,
		GrievanceRepo:   repos.Grievance,
		ReactionRepo:    repos.Reaction,
	}
}

// provisionalKRA returns a fallback DevelopmentScore when the compute call fails.
// Uses globalAverage blended to weight — identical to what Bayesian returns
// when n=0 — so a failed KRA doesn't silently zero out the total.
func (svc Score) ProvisionalDevelopmentKRA(weight uint) model.DevelopmentScore {
	return model.DevelopmentScore{
		Score:         (globalAverage / 100.0) * float64(weight),
		Weight:        weight,
		GlobalAverage: globalAverage,
		IsProvisional: true,
		Confidence:    0,
	}
}

func (svc Score) ProvisionalGrievanceKRA(weight uint) model.GrievanceScore {
	return model.GrievanceScore{
		Score:         (globalAverage / 100.0) * float64(weight),
		Weight:        weight,
		GlobalAverage: globalAverage,
		IsProvisional: true,
		Confidence:    0,
	}
}

func (svc Score) ProvisionalNewsPressKRA(weight uint) model.NewsPressScore {
	return model.NewsPressScore{
		Score:          (globalAverage / 100.0) * float64(weight),
		Weight:         weight,
		GlobalAverage:  globalAverage,
		IsProvisional:  true,
		Confidence:     0,
		LastComputedAt: time.Now(),
	}
}

func (svc Score) ProvisionalCALKRA(weight uint) model.CALRecordScore {
	return model.CALRecordScore{
		Score:         (globalAverage / 100.0) * float64(weight),
		Weight:        weight,
		GlobalAverage: globalAverage,
		IsProvisional: true,
		Confidence:    0,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────────────────

// applyKRAWeight scales a 0-100 score to the KRA's weight band.
// e.g. score=80, weight=40 → 32.0 (contributes 32 out of 40 possible points)
func applyKRAWeight(score float64, weight uint) float64 {
	return (score / 100.0) * float64(weight)
}

// confidenceScore returns how much the real data dominates the Bayesian prior.
// Result is 0–100: 0 = no reactions, ~50 = at bayesianPriorStrength reactions,
// approaching 100 asymptotically as reactions grow.
func confidenceScore(n float64) float64 {
	if (n + bayesianPriorStrength) == 0 {
		return 0
	}
	return math.Min(100, (n/(n+bayesianPriorStrength))*100)
}
