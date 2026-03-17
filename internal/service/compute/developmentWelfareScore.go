package compute

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/log"
)

// ─────────────────────────────────────────────────────────────────────────────
// DEVELOPMENT & WELFARE
// ─────────────────────────────────────────────────────────────────────────────

// ComputeDevelopmentScore calculates the Development & Welfare KRA score.
//
// Formula (transparent):
//
//	postAvg(p)   = mean(NormalizedScore) across all LOCKED reactions on post p
//	rawScore     = sum(postAvg) / totalPosts          ← denom = ALL posts, not just reacted
//	bayesScore   = (n×rawScore + B×globalAvg) / (n+B) ← n = total locked reaction count
//	kraScore     = (bayesScore / 100) × weight         ← scaled to KRA weight band
//
// Edge cases handled:
//   - No posts → score=0, provisional=true
//   - Posts exist but 0 locked reactions → pulls score toward globalAvg via Bayesian
//   - Negative NormalizedScore (downvotes) preserved — not clamped
//   - useDefault=true (new member, <6 months) → provisional flag set
func (svc Score) ComputeDevelopmentScore(weight uint, memberID int, computeYear int, useDefault bool) (*model.DevelopmentScore, error) {
	kraReport := model.DevelopmentScore{
		Score:          0,
		Weight:         weight,
		GlobalAverage:  globalAverage,
		DevelopmentIds: make([]string, 0),
		IsProvisional:  useDefault,
	}

	// Fetch only this year's posts — annual report card boundary
	developments, err := svc.DevelopmentRepo.GetByMemberIDAndYear(memberID, computeYear)
	if err != nil {
		return nil, err
	}
	totalPosts := len(developments)
	kraReport.TotalCount = uint(totalPosts)

	if totalPosts == 0 {
		// No posts at all — full provisional, score = Bayesian prior
		kraReport.Score = applyKRAWeight(globalAverage, weight)
		kraReport.Confidence = 0
		kraReport.IsProvisional = true
		kraReport.LastComputedAt = time.Now()
		return &kraReport, nil
	}

	// Goroutine pool — one per post, capped at NumCPU
	type postResult struct {
		postAvg       float64 // average NormalizedScore for this post
		reactionCount int     // number of locked reactions
		hasReactions  bool
	}

	results := make([]postResult, totalPosts)

	maxWorkers := runtime.NumCPU()
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, dev := range developments {
		wg.Add(1)
		sem <- struct{}{}
		idx, d := i, dev

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			// CRITICAL: only fetch locked reactions
			reactions, err := svc.ReactionRepo.GetLockedReactions(d.ID)
			if err != nil {
				log.Errorf("ComputeDevelopmentScore: failed to fetch locked reactions for post %v: %v", d.ID, err)
				return
			}
			if len(reactions) == 0 {
				// Post exists, no locked reactions yet.
				// Contributes 0 to numerator, counted in denominator — correct.
				return
			}

			var scoreSum float64
			for _, r := range reactions {
				// Preserve negative scores — downvote sentiment is real data.
				// NormalizedScore is already on 0–100 scale from SaveReaction.
				scoreSum += float64(r.NormalizedScore)
			}
			postAvg := scoreSum / float64(len(reactions))

			mu.Lock()
			results[idx] = postResult{
				postAvg:       postAvg,
				reactionCount: len(reactions),
				hasReactions:  true,
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Aggregate — denominator is ALL posts (totalPosts)
	var numeratorSum float64   // sum of per-post averages
	var totalReactionCount int // total locked reactions across all posts (used as Bayesian n)
	var postsWithReactions int
	for _, r := range results {
		if r.hasReactions {
			numeratorSum += r.postAvg
			totalReactionCount += r.reactionCount
			postsWithReactions++
		}
		// posts with 0 locked reactions add 0 to numerator — intentional
	}

	// rawScore: simple average of post with reaction averages, denom = all posts with reactions
	var rawScore float64
	if postsWithReactions > 0 {
		rawScore = numeratorSum / float64(postsWithReactions)
	}

	// Bayesian smoothing: blends rawScore toward globalAvg when n is small.
	// As totalReactionCount grows, weight on globalAvg shrinks naturally.
	n := float64(totalReactionCount)
	bayesScore := (n*rawScore + bayesianPriorStrength*globalAverage) / (n + bayesianPriorStrength)

	// Clamp to valid range after Bayesian (defensive)
	bayesScore = math.Max(0, math.Min(100, bayesScore))

	// Scale to KRA weight band — result is in [0, weight]
	kraReport.Score = applyKRAWeight(bayesScore, weight)

	// Confidence: how much the real data dominates vs the prior
	// 0% = no data at all, 100% = so many reactions prior is irrelevant
	kraReport.Confidence = confidenceScore(n)

	// Mark provisional if insufficient evidence
	if totalReactionCount < minEvidenceReactions {
		kraReport.IsProvisional = true
	}

	kraReport.LastComputedAt = time.Now()
	return &kraReport, nil
}
