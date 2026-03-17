package compute

import (
	"time"

	"github.com/shappy0/mantritva/internal/model"
)

// ─────────────────────────────────────────────────────────────────────────────
// CAL RECORDS (Criminal & Legal)
// ─────────────────────────────────────────────────────────────────────────────

// ComputeCALRecordScore calculates the Criminal & Legal Records KRA score.
//
// CAL records are factual/legal — they are NOT reaction-based.
// Score is derived from presence and nature of records:
//   - No records → full score (100)
//   - Each record reduces score based on severity
//
// TODO: implement when CAL record model and repo are ready.
// Returning provisional score (globalAverage) until then.
func (svc Score) ComputeCALRecordScore(weight uint, memberID int, computeYear int, useDefault bool) (*model.CALRecordScore, error) {
	kraReport := model.CALRecordScore{
		Score:          applyKRAWeight(globalAverage, weight),
		Weight:         weight,
		GlobalAverage:  globalAverage,
		CALRecordIds:   nil,
		IsProvisional:  true, // always provisional until implemented
		Confidence:     0,
		TotalCount:     0,
		LastComputedAt: time.Now(),
	}
	return &kraReport, nil
}
