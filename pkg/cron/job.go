package cron

import (
	"context"

	rcron "github.com/robfig/cron/v3"
)

type CronJob struct {
	Cron *rcron.Cron
}

func NewCronJob() *CronJob {
	cj := &CronJob{
		Cron: rcron.New(),
	}
	return cj
}

func (cj *CronJob) AddFunc(spec string, fn func()) error {
	_, err := cj.Cron.AddFunc(spec, fn)
	return err
}

func (cj *CronJob) Start() {
	cj.Cron.Start()
}

func (cj *CronJob) Stop() context.Context {
	return cj.Cron.Stop()
}
