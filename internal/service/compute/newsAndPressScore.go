package compute

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/log"
)

// ─────────────────────────────────────────────────────────────────────────────
// NEWS & PRESS STATEMENTS
// ─────────────────────────────────────────────────────────────────────────────

// ComputeNewsPressScore calculates the News & Press Statements KRA score.
//
// Same formula as Development. Bayesian applied ONCE at aggregate level
// (not per-post — that double-smooths and understates variance).
func (svc Score) ComputeNewsPressScore(weight uint, memberID int, computeYear int, useDefault bool) (*model.NewsPressScore, error) {
	kraReport := model.NewsPressScore{
		Score:         0,
		Weight:        weight,
		GlobalAverage: globalAverage,
		NewsPressIds:  make([]string, 0),
		IsProvisional: useDefault,
	}

	pagination := dto.Pagination{Offset: 0, Limit: maxPostsPerFetch}
	newsPresses, err := svc.NewsPressRepo.GetByMemberIDAndYear(pagination, memberID, computeYear)
	if err != nil {
		return nil, err
	}
	totalPosts := len(*newsPresses)
	kraReport.TotalCount = uint(totalPosts)

	if totalPosts == 0 {
		kraReport.Score = applyKRAWeight(globalAverage, weight)
		kraReport.Confidence = 0
		kraReport.IsProvisional = true
		kraReport.LastComputedAt = time.Now()
		return &kraReport, nil
	}

	type postResult struct {
		postAvg       float64
		reactionCount int
		hasReactions  bool
	}
	results := make([]postResult, totalPosts)

	maxWorkers := runtime.NumCPU()
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, post := range *newsPresses {
		wg.Add(1)
		sem <- struct{}{}
		idx, p := i, post

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			reactions, err := svc.ReactionRepo.GetLockedReactions(p.ID)
			if err != nil {
				log.Errorf("ComputeNewsPressScore: failed to fetch locked reactions for post %v: %v", p.ID, err)
				return
			}
			if len(reactions) == 0 {
				return
			}

			var scoreSum float64
			for _, r := range reactions {
				scoreSum += float64(r.NormalizedScore)
			}

			mu.Lock()
			results[idx] = postResult{
				postAvg:       scoreSum / float64(len(reactions)),
				reactionCount: len(reactions),
				hasReactions:  true,
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	var numeratorSum float64
	var totalReactionCount int
	var postsWithReactions int
	for _, r := range results {
		if r.hasReactions {
			numeratorSum += r.postAvg
			totalReactionCount += r.reactionCount
			postsWithReactions++
		}
	}
	var rawScore float64
	if postsWithReactions > 0 {
		rawScore = numeratorSum / float64(postsWithReactions)
	}

	n := float64(totalReactionCount)
	bayesScore := (n*rawScore + bayesianPriorStrength*globalAverage) / (n + bayesianPriorStrength)
	bayesScore = math.Max(0, math.Min(100, bayesScore))

	kraReport.Score = applyKRAWeight(bayesScore, weight)
	kraReport.Confidence = confidenceScore(n)
	if totalReactionCount < minEvidenceReactions {
		kraReport.IsProvisional = true
	}
	kraReport.LastComputedAt = time.Now()
	return &kraReport, nil
}
