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

func MakeTicker(year int, month time.Month, day, hour, min, sec int) *Alarm {
	now := time.Now()
	endTime := targetTime(now, year, month, day, hour, min, sec)
	duration := endTime.Sub(now)

	alarm := new(Alarm)
	alarm.Ticker = time.NewTicker(duration)
	alarm.FinishesAt = endTime
	alarm.Alert = make(chan struct{})

	return alarm
}

func (a *Alarm) UpdateTicker(year int, month time.Month, day, hour, min, sec int) *Alarm {
	close(a.Alert)

	now := time.Now()
	endTime := targetTime(now, year, month, day, hour, min, sec)
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
