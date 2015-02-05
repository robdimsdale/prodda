package timer

import (
	"time"
)

type Alarm struct {
	Ticker     *time.Ticker
	FinishesAt time.Time
}

func MakeTimer(year int, month time.Month, day, hour, min, sec int) *Alarm {
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
