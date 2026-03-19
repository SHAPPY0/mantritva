package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service/compute"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type ReportService interface {
	GenerateReport(int) (*model.Score, error)
	GetReport(int) (dto.MemberReport, error)
	GetMemberReports(int, int, int) ([]dto.MemberReport, error)
	GeneratReportCron()
	DownloadReport(string, string) ([]byte, error)
}

type reportService struct {
	ReportRepo      repository.ReportRepository
	GrievanceRepo   repository.GrievanceRepository
	DevelopmentRepo repository.DevelopmentWelfareRepository
	NewsPressRepo   repository.NewsPressPostRepository
	PartyRepo       repository.PartyRepository
	ReactionRepo    repository.ReactionRepository
	ScoreSvc        compute.Score
}

func NewReportService(repos *repository.Repos) ReportService {
	return &reportService{
		ReportRepo:      repos.Report,
		GrievanceRepo:   repos.Grievance,
		DevelopmentRepo: repos.DevelopmentWelfare,
		NewsPressRepo:   repos.NewsPressPost,
		PartyRepo:       repos.Party,
		ReactionRepo:    repos.Reaction,
		ScoreSvc:        compute.NewScore(repos),
	}
}

func (s *reportService) GenerateReport(memberID int) (*model.Score, error) {
	computeYear := time.Now().Year()
	log.Infof("GenerateReport: starting for memberID=%d year=%d", memberID, computeYear)

	// ── Fetch KRA weights from DB ─────────────────────────────────────────
	var kraWeights dto.KRAWeights
	kras, err := s.ReportRepo.GetReportKRA()
	if err != nil {
		return nil, err
	}
	for _, kra := range kras {
		switch kra.KRAName {
		case string(utils.KRA_DEVELOPMENTS):
			kraWeights.Developments = kra.Weight
		case string(utils.KRA_GRIEVANCES):
			kraWeights.Grievances = kra.Weight
		case string(utils.KRA_CALRECORDS):
			kraWeights.CALRecords = kra.Weight
		case string(utils.KRA_NEWSPRESSES):
			kraWeights.NewsPresses = kra.Weight
		}
	}

	// ── Determine if member is new (< 6 months political career) ─────────
	// New members don't have enough data. Their KRA scores will be blended
	// toward globalAverage with is_provisional=true until data accumulates.
	useDefault := s.isMemberNew(memberID)

	// ── Compute each KRA concurrently ─────────────────────────────────────
	// Each Compute call is independent — run in parallel.
	type kraResult struct {
		developments *model.DevelopmentScore
		grievances   *model.GrievanceScore
		newsPresses  *model.NewsPressScore
		calRecords   *model.CALRecordScore
	}
	var (
		res  kraResult
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []string
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		r, err := s.ScoreSvc.ComputeDevelopmentScore(kraWeights.Developments, memberID, computeYear, useDefault)
		if err != nil {
			log.Errorf("GenerateReport: development score failed for member %d: %v", memberID, err)
			mu.Lock()
			errs = append(errs, "development")
			mu.Unlock()
			return
		}
		mu.Lock()
		res.developments = r
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		r, err := s.ScoreSvc.ComputeGrievanceScore(kraWeights.Grievances, memberID, computeYear, useDefault)
		if err != nil {
			log.Errorf("GenerateReport: grievance score failed for member %d: %v", memberID, err)
			mu.Lock()
			errs = append(errs, "grievance")
			mu.Unlock()
			return
		}
		mu.Lock()
		res.grievances = r
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		r, err := s.ScoreSvc.ComputeNewsPressScore(kraWeights.NewsPresses, memberID, computeYear, useDefault)
		if err != nil {
			log.Errorf("GenerateReport: newspress score failed for member %d: %v", memberID, err)
			mu.Lock()
			errs = append(errs, "newspress")
			mu.Unlock()
			return
		}
		mu.Lock()
		res.newsPresses = r
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		r, err := s.ScoreSvc.ComputeCALRecordScore(kraWeights.CALRecords, memberID, computeYear, useDefault)
		if err != nil {
			log.Errorf("GenerateReport: CAL record score failed for member %d: %v", memberID, err)
			mu.Lock()
			errs = append(errs, "cal")
			mu.Unlock()
			return
		}
		mu.Lock()
		res.calRecords = r
		mu.Unlock()
	}()

	wg.Wait()

	// ── Assemble Score document ───────────────────────────────────────────
	scoreData := model.Score{
		MemberID:        uint(memberID),
		MaxScore:        100,
		Year:            computeYear,
		AuditFieldsBSON: model.InsertBSONAuditFields("SYSTEM"),
	}

	var totalScore float64
	var totalConfidenceWeighted float64
	var totalWeightUsed float64

	if res.developments != nil {
		scoreData.Developments = *res.developments
		totalScore += res.developments.Score
		totalConfidenceWeighted += res.developments.Confidence * float64(kraWeights.Developments)
		totalWeightUsed += float64(kraWeights.Developments)
	} else {
		// KRA failed — use provisional fallback so total score is not silently wrong
		scoreData.Developments = s.ScoreSvc.ProvisionalDevelopmentKRA(kraWeights.Developments)
		totalScore += scoreData.Developments.Score
		totalWeightUsed += float64(kraWeights.Developments)
	}

	if res.grievances != nil {
		scoreData.Grievances = *res.grievances
		totalScore += res.grievances.Score
		totalConfidenceWeighted += res.grievances.Confidence * float64(kraWeights.Grievances)
		totalWeightUsed += float64(kraWeights.Grievances)
	} else {
		scoreData.Grievances = s.ScoreSvc.ProvisionalGrievanceKRA(kraWeights.Grievances)
		totalScore += scoreData.Grievances.Score
		totalWeightUsed += float64(kraWeights.Grievances)
	}

	if res.newsPresses != nil {
		scoreData.NewsPresses = *res.newsPresses
		totalScore += res.newsPresses.Score
		totalConfidenceWeighted += res.newsPresses.Confidence * float64(kraWeights.NewsPresses)
		totalWeightUsed += float64(kraWeights.NewsPresses)
	} else {
		scoreData.NewsPresses = s.ScoreSvc.ProvisionalNewsPressKRA(kraWeights.NewsPresses)
		totalScore += scoreData.NewsPresses.Score
		totalWeightUsed += float64(kraWeights.NewsPresses)
	}

	if res.calRecords != nil {
		scoreData.CALRecords = *res.calRecords
		totalScore += res.calRecords.Score
		totalConfidenceWeighted += res.calRecords.Confidence * float64(kraWeights.CALRecords)
		totalWeightUsed += float64(kraWeights.CALRecords)
	} else {
		scoreData.CALRecords = s.ScoreSvc.ProvisionalCALKRA(kraWeights.CALRecords)
		totalScore += scoreData.CALRecords.Score
		totalWeightUsed += float64(kraWeights.CALRecords)
	}

	// totalScore is the direct sum of KRA scores, each already in their weight band.
	// e.g. dev=32 + grievance=21 + news=10 + cal=8 = 71 out of 100.
	scoreData.TotalScore = totalScore
	if totalWeightUsed > 0 {
		scoreData.Confidence = totalConfidenceWeighted / totalWeightUsed
	}

	// Mark overall as provisional if any KRA is provisional
	scoreData.IsProvisional = scoreData.Developments.IsProvisional ||
		scoreData.Grievances.IsProvisional ||
		scoreData.NewsPresses.IsProvisional ||
		scoreData.CALRecords.IsProvisional

	log.Infof("***ScoreData*** %v", scoreData)
	// ── Upsert — never insert duplicates ──────────────────────────────────
	// Use Upsert(filter={member_id, year}, update=$set{...}) in MongoDB.
	// This replaces the old Create() which inserted a new doc every run.
	if err := s.ReportRepo.CreateScore(&scoreData); err != nil {
		return nil, err
	}

	log.Infof("GenerateReport: completed memberID=%d year=%d totalScore=%.2f provisional=%v failedKRAs=%v",
		memberID, computeYear, scoreData.TotalScore, scoreData.IsProvisional, errs)
	return &scoreData, nil
}

func (svc *reportService) GeneratReportCron() {
	computeDateTime := time.Now().UTC()
	log.Infof("GenerateReportCron: Starting run at %s", computeDateTime.Format(time.RFC3339))

	//Step1: fetch reactions that have passed their lock tim
	lockedReactions, err := svc.ReactionRepo.GetAndLockMaturedReactions(computeDateTime)
	if err != nil {
		log.Errorf("GenerateReportCron: Failed to fetch matured reactions: %v", err)
		return
	}
	if len(lockedReactions) == 0 {
		log.Infof("GenerateReportCron: No matured reactions - nothing to do")
		return
	}
	log.Infof("GenerateReportCronn: %d reactions newly locked", len(lockedReactions))

	//Step 2: Resolve memberID for each reaction and deduplicate
	memberIDs := svc.resolveMemberIDs(lockedReactions)
	if len(memberIDs) == 0 {
		log.Warnf("GenerateReportCron: Reaction found but no memberIDs resolved - check post linkage")
		return
	}
	log.Infof("GenerateReportCron: %d distinct members to recompute", len(memberIDs))

	//Step 3: Recompute scores - bounded concurreny (max 5 members at a time)
	const maxConcurrentMembers = 5
	sem := make(chan struct{}, maxConcurrentMembers)
	var wg sync.WaitGroup

	for _, memberID := range memberIDs {
		wg.Add(1)
		sem <- struct{}{}
		mID := memberID
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			if _, err := svc.GenerateReport(int(mID)); err != nil {
				log.Errorf("GenerateReportCron: GenerateReport failed for member %d: %v", mID, err)
			}
		}()
	}
	wg.Wait()
	log.Infof("GenerateReportCron: Completed Run - Processed %d members", len(memberIDs))
}

func (svc *reportService) resolveMemberIDs(reactions []model.Reaction) []uint {
	seen := make(map[uint]struct{})
	var result []uint
	for _, reaction := range reactions {
		postID := reaction.PostID
		memberID := uint(0)
		if developemnt, err := svc.DevelopmentRepo.GetByID(postID); err == nil && developemnt != nil && developemnt.MemberID != 0 {
			memberID = developemnt.MemberID
		} else if grievance, err := svc.GrievanceRepo.GetByID(postID); err == nil && grievance != nil {
			switch {
			case grievance.MemberOfParliamentID != 0:
				memberID = grievance.MemberOfParliamentID
			case grievance.StateLeaderID != 0:
				memberID = grievance.StateLeaderID
			case grievance.CabinetMinisterID != 0:
				memberID = grievance.CabinetMinisterID
			}
		} else if np, err := svc.NewsPressRepo.GetByID(postID); err == nil && np != nil && np.MemberID != 0 {
			memberID = np.MemberID
		}

		if memberID == 0 {
			log.Warnf("resolveMemberIDs: Cound not resolve memberID for post %v", postID)
			continue
		}
		if _, exists := seen[memberID]; !exists {
			seen[memberID] = struct{}{}
			result = append(result, memberID)
		}
	}
	return result
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

func (s *reportService) DownloadReport(url, divSelector string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Infof))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var divHTML string
	var pdfBuf []byte
	var cssContent string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),

		// ── Fix 1: Evaluate with awaitPromise ────────────────────────────
		chromedp.ActionFunc(func(ctx context.Context) error {
			var res interface{}
			return chromedp.Evaluate(`
				new Promise((resolve) => {
					if (document.readyState === 'complete') {
						setTimeout(resolve, 3000)
					} else {
						window.addEventListener('load', () => setTimeout(resolve, 3000))
					}
				})
			`, &res, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
				return p.WithAwaitPromise(true) // ← critical, without this Promise is not awaited
			}).Do(ctx)
		}),

		// ── Fix 2: Poll for div instead of WaitVisible ───────────────────
		// WaitVisible can timeout if selector doesn't exist
		// This polls every 500ms for up to 60s
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Poll(fmt.Sprintf(`document.querySelector('#%s') !== null`, divSelector), nil,
				chromedp.WithPollingInterval(500*time.Millisecond),
				chromedp.WithPollingTimeout(60*time.Second),
			).Do(ctx)
		}),
		chromedp.WaitVisible("#"+divSelector, chromedp.ByQuery),

		// ── extract all styles from the live page ──────────────────────
		// this grabs every <style> tag and every loaded <link> stylesheet
		// already parsed and applied by the browser — no external requests needed
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Evaluate(`
                (function() {
                    let css = '';
                    // grab all stylesheets already loaded by the browser
                    for (const sheet of document.styleSheets) {
                        try {
                            for (const rule of sheet.cssRules) {
                                css += rule.cssText + '\n';
                            }
                        } catch(e) {
                            // cross-origin sheets may throw — skip
                        }
                    }
                    return css;
                })()
            `, &cssContent).Do(ctx)
		}),

		// ── After capturing divHTML, convert all image srcs to base64 ─────────
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Evaluate(`
				(async function() {
					const imgs = document.querySelectorAll('#`+divSelector+` img');
					const toBase64 = (url) => fetch(url)
						.then(r => r.blob())
						.then(blob => new Promise((resolve, reject) => {
							const reader = new FileReader();
							reader.onloadend = () => resolve(reader.result);
							reader.onerror = reject;
							reader.readAsDataURL(blob);
						}));

					for (const img of imgs) {
						try {
							const src = img.src; // always absolute in browser
							if (src && !src.startsWith('data:')) {
								img.src = await toBase64(src);
							}
						} catch(e) {
							console.warn('failed to convert image:', img.src, e);
						}
					}
				})()
			`, nil, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
				return p.WithAwaitPromise(true) // wait for all images to convert
			}).Do(ctx)
		}),

		chromedp.OuterHTML("#"+divSelector, &divHTML, chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("failed to capture div: %w", err)
	}

	if divHTML == "" {
		return nil, fmt.Errorf("div #%s not found or empty", divSelector)
	}

	fullHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				* {
					-webkit-print-color-adjust: exact !important;
					print-color-adjust: exact !important;
					margin: 0;
					padding: 0;
					box-sizing: border-box;
				}
				body {
					font-family: sans-serif;
				}
				a::after {
					display: none !important;
					content: none !important;
				}
				a {
					text-decoration: none; /* optional — remove underline too */
				}
				%s
			</style>
			
		</head>
		<body>
			%s
		</body>
		</html>
	`, cssContent, divHTML)
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, fullHTML).Do(ctx)
		}),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(1*time.Second), // let CSS load
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0.1).
				WithMarginBottom(0.1).
				WithMarginLeft(0.1).
				WithMarginRight(0.1).
				Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("failed to generate pdf: %w", err)
	}

	return pdfBuf, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────────────────

// isMemberNew returns true if the member started their political career
// less than 6 months ago — i.e. insufficient historical data to produce
// a meaningful score from evidence alone.
func (s *reportService) isMemberNew(memberID int) bool {
	startDate, err := s.PartyRepo.GetPoliticalStartDate(uint(memberID))
	if err != nil {
		log.Warnf("isMemberNew: could not fetch political start date for member %d: %v — defaulting to false", memberID, err)
		return false
	}
	parsed, err := time.Parse("2006-01-02", startDate.Format("2006-01-02"))
	if err != nil {
		return false
	}
	// < 6 months = new member, use provisional scores
	return time.Now().Before(parsed.AddDate(0, 6, 0))
}
