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
	targetTime := time.Date(year, month, day, hour, min, sec, 0, time.Local)
	if !targetTime.After(now) {
		targetTime = targetTime.Add(time.Duration(24) * time.Hour)
	}
	duration := targetTime.Sub(now)

	alarm := new(Alarm)
	alarm.Ticker = time.NewTicker(duration)
	alarm.FinishesAt = targetTime

	return alarm
}
