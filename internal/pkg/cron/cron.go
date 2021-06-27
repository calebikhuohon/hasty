package cron

import (
	"time"
)

const (
	IntervalPeriod     = 2 * time.Hour
	HourToTick     int = 23
	MinuteToTick   int = 21
	SecondToTick   int = 03
)

type JobTicker struct {
	T *time.Timer
}

func GetNextTickDuration() time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), HourToTick, MinuteToTick, SecondToTick, 0, time.Local)
	if nextTick.Before(now) {
		nextTick = nextTick.Add(IntervalPeriod)
	}
	return nextTick.Sub(time.Now())
}

func NewJobTicker() JobTicker {
	return JobTicker{time.NewTimer(GetNextTickDuration())}
}

func (jt JobTicker) UpdateJobTicker() {
	jt.T.Reset(GetNextTickDuration())
}
