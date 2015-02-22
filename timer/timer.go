package timer

import (
	"errors"
	"time"
)

type Alarm struct {
	running    bool
	Ticker     *time.Ticker
	FinishesAt time.Time
	Alert      chan struct{}
	task       Task
}

type Task interface {
	Run() error
}

// NewAlarm creates an alarm
// t must be after the current time.
func NewAlarm(t time.Time, task Task) (*Alarm, error) {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return nil, errors.New("Time must not be in the past")
	}

	duration := t.Sub(currentTime)

	return &Alarm{
		FinishesAt: t,
		Alert:      make(chan struct{}),
		Ticker:     time.NewTicker(duration),
		task:       task,
	}, nil
}

func (a *Alarm) UpdateAlarm(t time.Time) error {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return errors.New("Time must not be in the past")
	}

	if a.running {
		close(a.Alert)
	}

	a.FinishesAt = t
	duration := t.Sub(currentTime)
	a.Ticker = time.NewTicker(duration)
	a.Alert = make(chan struct{})

	return nil
}

func (a *Alarm) RunOnDing() error {
	a.running = true
	select {
	case <-a.Ticker.C:
		err := a.task.Run()
		if err != nil {
			return err
		}
		a.running = false
	case <-a.Alert:
		a.Ticker.Stop()
		a.running = false
	}
	return nil
}
