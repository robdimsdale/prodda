package timer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mfine30/prodda/client"
)

type Alarm struct {
	running    bool
	Ticker     *time.Ticker
	FinishesAt time.Time
	CancelChan chan struct{}
	task       Task
}

type Task interface {
	Run() error
}

type TravisTask struct {
	client  *client.Travis
	token   string
	buildID uint
}

func NewTravisTask(token string, buildID uint) *TravisTask {
	return &TravisTask{
		client:  client.NewTravisClient("https://api.travis-ci.org"),
		token:   token,
		buildID: buildID,
	}
}

func (t TravisTask) Run() error {
	fmt.Printf("Travis task running\n")

	resp, err := t.client.TriggerBuild(t.token, t.buildID)
	if err != nil {
		return err
	}
	log.Printf("response: %+v\n", resp)
	return nil
}

// NewAlarm creates an alarm
// t must be after the current time.
// task must be non-nil
func NewAlarm(t time.Time, task Task) (*Alarm, error) {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return nil, errors.New("Time must not be in the past")
	}

	if task == nil {
		return nil, errors.New("Task must not be nil.")
	}

	duration := t.Sub(currentTime)

	return &Alarm{
		FinishesAt: t,
		CancelChan: make(chan struct{}),
		Ticker:     time.NewTicker(duration),
		task:       task,
	}, nil
}

// Update will change the time the alarm will finish at.
// If the alarm is currently running it will be canceled.
// An error will be thrown if time is not in the future
func (a *Alarm) Update(t time.Time) error {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return errors.New("Time must not be in the past")
	}

	if a.running {
		a.Cancel()
	}

	a.FinishesAt = t
	duration := t.Sub(currentTime)
	a.Ticker = time.NewTicker(duration)
	a.CancelChan = make(chan struct{})

	return nil
}

// Cancel will cancel the alarm if it is running
// It will return an error if the alarm is not running.
func (a *Alarm) Cancel() error {
	if !a.running {
		return errors.New("Alarm not running")
	}

	close(a.CancelChan)
	return nil
}

// Start will block until either the timer goes off or the Alarm is canceled
func (a *Alarm) Start() error {
	a.running = true
	select {
	case <-a.Ticker.C:
		fmt.Printf("Alarm time has gone off\n")
		err := a.task.Run()
		if err != nil {
			fmt.Printf("Error running task: %v\n", err)
			return err
		}
		a.running = false
	case <-a.CancelChan:
		fmt.Printf("Alarm canceled\n")
		a.Ticker.Stop()
		a.running = false
	}
	return nil
}
