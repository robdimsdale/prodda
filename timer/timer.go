package timer

import (
	"errors"
	"time"
)

type Alarm struct {
	started    bool
	Ticker     *time.Ticker
	FinishesAt time.Time
	Alert      chan struct{}
}

type Task interface {
	Run() error
}

// NewAlarm creates an alarm
// t must be after the current time.
func NewAlarm(t time.Time) (*Alarm, error) {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return nil, errors.New("Time must not be in the past")
	}

	duration := t.Sub(currentTime)

	return &Alarm{
		FinishesAt: t,
		Alert:      make(chan struct{}),
		Ticker:     time.NewTicker(duration),
	}, nil
}

func (a *Alarm) UpdateAlarm(t time.Time) error {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return errors.New("Time must not be in the past")
	}

	if a.started {
		close(a.Alert)
	}

	a.FinishesAt = t
	duration := t.Sub(currentTime)
	a.Ticker = time.NewTicker(duration)
	a.Alert = make(chan struct{})
	a.started = false //TODO: backfill test for this

	return nil
}

func (a *Alarm) RunOnDing(task Task) error {
	a.started = true
	select {
	case <-a.Ticker.C:
		err := task.Run()
		if err != nil {
			return err
		}
		a.started = false
	case <-a.Alert:
		a.Ticker.Stop()
		a.started = false
	}
	return nil
}
