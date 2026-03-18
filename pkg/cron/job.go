package cron

import (
	"context"

	rcron "github.com/robfig/cron/v3"
)

type CronJob struct {
	Cron *rcron.Cron
}

func NewCronJob(fn func()) *CronJob {
	cj := &CronJob{
		Cron: rcron.New(),
	}
	// cj.Cron.AddFunc("@every 5s", fn)
	cj.Cron.AddFunc("@every 15m", fn)
	return cj
}

func (cj *CronJob) Start() {
	cj.Cron.Start()
}

func (cj *CronJob) Stop() context.Context {
	return cj.Cron.Stop()
}
