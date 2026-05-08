package service

import (
	"sort"
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
)

const (
	TOP_SCORER_LIMIT = 5
)

type AnalyticStatsService interface {
	ComputeTopScorer(*model.Score) error
	GetTopScorer(int) (*model.AnalyticStat, error)
}

type analyticStatService struct {
	Repository repository.AnalyticStatRepository
}

func NewAnalyticStatService(repo *repository.Repos) AnalyticStatsService {
	return &analyticStatService{
		Repository: repo.AnalyticStat,
	}
}

func (svc *analyticStatService) ComputeTopScorer(report *model.Score) error {
	STAT_TYPE := "TOP_SCORER"
	currentDateTime := time.Now()
	currentYear := currentDateTime.Year()
	currentMonth := currentDateTime.Month()
	existingStat, err := svc.Repository.GetByYear(STAT_TYPE, currentYear)
	if err != nil {
	}
	rankData := map[string]any{
		"rank":       1,
		"score":      report.TotalScore,
		"confidence": report.Confidence,
		"member_id":  report.MemberID,
	}

	if existingStat == nil {
		stat := &model.AnalyticStat{
			StatType: STAT_TYPE,
			Scope:    "Global",
			Year:     currentYear,
			Month:    int(currentMonth),
			Data:     []map[string]any{rankData},
			MetaData: map[string]any{
				"rank_limit":  TOP_SCORER_LIMIT,
				"computed_at": time.Now(),
			},
			AuditFieldsBSON: model.InsertBSONAuditFields("SYSTEM"),
		}
		if _, err := svc.Repository.Create(stat); err != nil {
			log.Errorf("ComputeTopScorer: Error storing analyticStat %v", err)
		}
		return nil
	}
	inserted := false
	for i, entry := range existingStat.Data {
		// Check if memberID already exists
		if mID, ok := entry["member_id"]; ok {
			if memberID, ok := mID.(uint); ok && memberID == report.MemberID {
				if s, ok := toFloat(entry["score"]); ok && report.TotalScore > s {
					existingStat.Data[i] = rankData
				}
				inserted = true
				break
			}
		}
		// Check score if memberID doesn't match
		if s, ok := toFloat(entry["score"]); ok {
			if report.TotalScore > s {
				newSlice := make([]map[string]any, 0, len(existingStat.Data)+1)
				newSlice = append(newSlice, existingStat.Data[:i]...)
				newSlice = append(newSlice, rankData)
				newSlice = append(newSlice, existingStat.Data[i:]...)
				existingStat.Data = newSlice
				inserted = true
				break
			}
		}
	}

	// if not inserted and there's room, append at end
	if !inserted && len(existingStat.Data) < TOP_SCORER_LIMIT {
		existingStat.Data = append(existingStat.Data, rankData)
	}

	// ensure length limit
	if len(existingStat.Data) > TOP_SCORER_LIMIT {
		existingStat.Data = existingStat.Data[:TOP_SCORER_LIMIT]
	}

	// recompute ranks based on position (1-based)
	sort.Slice(existingStat.Data, func(i, j int) bool {
		scoreI := existingStat.Data[i]["score"].(float64)
		scoreJ := existingStat.Data[j]["score"].(float64)
		return scoreI > scoreJ // descending
	})

	for i := range existingStat.Data {
		existingStat.Data[i]["rank"] = i + 1
	}

	if _, err := svc.Repository.Update(existingStat); err != nil {
		return err
	}

	return nil
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	default:
		return 0, false
	}
}

func (svc *analyticStatService) GetTopScorer(year int) (*model.AnalyticStat, error) {
	scorer, err := svc.Repository.GetByYear("TOP_SCORER", year)
	if err != nil {
		return nil, err
	}
	return scorer, nil
}
