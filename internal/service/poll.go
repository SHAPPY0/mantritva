package service

import (
	"fmt"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

const (
	PollYesNo       = "YES_NO"
	PollFourOptions = "FOUR_OPTIONS"
	PollActive      = "ACTIVE"
	PollClosed      = "CLOSED"
)

type PollService interface {
	Create(int, dto.PollCreateReq) (*dto.PollRes, error)
	Get(int, primitive.ObjectID) (*dto.PollRes, error)
	Vote(int, primitive.ObjectID, int) error
	CloseExpired() error
}
type pollService struct {
	Repo repository.PollRepository
	Feed repository.FeedRepository
}

func NewPollService(r *repository.Repos) PollService { return &pollService{Repo: r.Poll, Feed: r.Feed} }
func (s *pollService) Create(user int, req dto.PollCreateReq) (*dto.PollRes, error) {
	q := strings.TrimSpace(req.Question)
	if q == "" || len([]rune(q)) > 500 {
		return nil, fmt.Errorf("question is required and must be at most 500 characters")
	}
	typ := strings.ToUpper(strings.TrimSpace(req.PollType))
	if typ != PollYesNo && typ != PollFourOptions {
		return nil, fmt.Errorf("invalid poll type")
	}
	var texts []string
	if typ == PollYesNo {
		texts = []string{"Yes", "No"}
	} else {
		if len(req.Options) != 4 {
			return nil, fmt.Errorf("FOUR_OPTIONS requires exactly four options")
		}
		seen := map[string]bool{}
		for _, raw := range req.Options {
			t := strings.TrimSpace(raw)
			key := strings.ToLower(t)
			if t == "" || len([]rune(t)) > 150 || seen[key] {
				return nil, fmt.Errorf("options must be non-empty, unique, and at most 150 characters")
			}
			seen[key] = true
			texts = append(texts, t)
		}
	}
	d := req.DurationDays
	if d == 0 {
		d = 3
	}
	if d != 1 && d != 3 && d != 7 {
		return nil, fmt.Errorf("duration must be 1, 3, or 7 days")
	}
	now := time.Now()
	p := &model.Poll{ID: primitive.NewObjectID(), Question: q, PollType: typ, DurationDays: d, StartsAt: now, EndsAt: now.AddDate(0, 0, d), Status: PollActive, AuditFieldsBSON: model.InsertBSONAuditFields(fmt.Sprintf("%d", user))}
	for i, t := range texts {
		p.Options = append(p.Options, model.PollOption{ID: i + 1, Text: t})
	}
	if err := s.Repo.Create(p); err != nil {
		return nil, err
	}
	f := &model.Feed{ID: primitive.NewObjectID(), Type: string(utils.POLL), PostID: p.ID, Status: utils.Published, AuditFieldsBSON: model.InsertBSONAuditFields(fmt.Sprintf("%d", user))}
	if _, err := s.Feed.Create(f); err != nil {
		return nil, err
	}
	p.FeedID = f.ID
	return s.response(p, false, nil), nil
}
func (s *pollService) Get(user int, id primitive.ObjectID) (*dto.PollRes, error) {
	p, e := s.Repo.GetByID(id)
	if e != nil {
		return nil, e
	}
	if p.Status == PollActive && time.Now().After(p.EndsAt) {
		p.Status = PollClosed
	}
	v, e := s.Repo.HasVoted(id, user)
	if e != nil {
		return nil, e
	}
	var selected *int
	if v != nil {
		selected = &v.OptionID
	}
	return s.response(p, v != nil, selected), nil
}
func (s *pollService) Vote(user int, id primitive.ObjectID, opt int) error {
	p, e := s.Repo.GetByID(id)
	if e != nil {
		return e
	}
	now := time.Now()
	if p.Status != PollActive || now.Before(p.StartsAt) || !now.Before(p.EndsAt) {
		return fmt.Errorf("poll is closed or not started")
	}
	valid := false
	for _, o := range p.Options {
		if o.ID == opt {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("option does not belong to poll")
	}
	v, e := s.Repo.HasVoted(id, user)
	if e != nil {
		return e
	}
	if v != nil {
		return fmt.Errorf("user has already voted")
	}
	return s.Repo.Vote(&model.PollVote{ID: primitive.NewObjectID(), PollID: id, OptionID: opt, UserID: user, CreatedAt: now})
}
func (s *pollService) CloseExpired() error { return s.Repo.CloseExpired() }
func (s *pollService) response(p *model.Poll, voted bool, selected *int) *dto.PollRes {
	r := &dto.PollRes{ID: p.ID, Question: p.Question, PollType: p.PollType, StartsAt: p.StartsAt, EndsAt: p.EndsAt, Status: p.Status, HasVoted: voted, SelectedOptionID: selected, Options: dto.PollOptions(p)}
	counts, _ := s.Repo.Results(p.ID)
	var total int64
	for _, n := range counts {
		total += n
	}
	for i := range r.Options {
		r.Options[i].Votes = counts[r.Options[i].ID]
		if total > 0 {
			r.Options[i].Percentage = float64(r.Options[i].Votes) * 100 / float64(total)
		}
	}
	return r
}
