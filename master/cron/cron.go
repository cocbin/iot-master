package cron

import (
	"github.com/go-co-op/gocron"
	"time"
)

var Scheduler *gocron.Scheduler

func init() {
	Scheduler = gocron.NewScheduler(time.UTC)
}

//是否是使用单一协程？？？ 是则要改成协程池？？？

func Schedule(crontab string, fn func()) (*Job, error) {
	job, err := Scheduler.Cron(crontab).Do(fn)
	if err != nil {
		return nil, err
	}
	return &Job{job: job}, nil
}

func Interval(interval int, fn func()) (*Job, error) {
	job, err := Scheduler.Every(interval).Milliseconds().Do(fn)
	if err != nil {
		return nil, err
	}
	return &Job{job: job}, nil
}

func Clock(hours int, minutes int, fn func()) (*Job, error) {
	job, err := Scheduler.At(hours).Hours().At(minutes).Minutes().Do(fn)
	if err != nil {
		return nil, err
	}
	return &Job{job: job}, nil
}

type Job struct {
	job *gocron.Job
}

func (j *Job) Cancel() {
	Scheduler.Remove(j.job)
}
