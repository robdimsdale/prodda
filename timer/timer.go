package timer

import (
	"time"
)

type Alarm struct {
	Ticker     *time.Ticker
	FinishesAt time.Time
	Alert      chan struct{}
}

type Task interface {
	Run() error
}

func NewTicker(year int, month time.Month, day, hour, min, sec int) *Alarm {
	now := time.Now()
	endTime := targetTime(now, year, month, day, hour, min, sec)
	alarm := new(Alarm)
	return alarm.makeTicker(now, endTime)
}

func (a *Alarm) UpdateTicker(year int, month time.Month, day, hour, min, sec int) *Alarm {
	close(a.Alert)

	now := time.Now()
	endTime := targetTime(now, year, month, day, hour, min, sec)
	a.makeTicker(now, endTime)

	return a
}

func (alarm *Alarm) RunOnDing(task Task) error {
	select {
	case <-alarm.Ticker.C:
		err := task.Run()
		if err != nil {
			return err
		}
	case <-alarm.Alert:
		alarm.Ticker.Stop()
	}
	return nil
}

func (a *Alarm) makeTicker(now, endTime time.Time) *Alarm {
	duration := endTime.Sub(now)

	a.Ticker = time.NewTicker(duration)
	a.FinishesAt = endTime
	a.Alert = make(chan struct{})

	return a
}

func targetTime(now time.Time, year int, month time.Month, day, hour, min, sec int) time.Time {
	endTime := time.Date(year, month, day, hour, min, sec, 0, time.Local)
	if !endTime.After(now) {
		endTime = endTime.Add(time.Duration(24) * time.Hour)
	}
	return endTime
}
