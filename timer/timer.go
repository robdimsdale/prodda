package timer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mfine30/prodda/client"
)

const (
	MinimumFrequency = time.Duration(1 * time.Minute)
)

type Alarm struct {
	running    bool
	NextTime   time.Time
	cancelChan chan struct{}
	task       Task
	frequency  time.Duration
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
// An error will be thrown if time is not in the future
// An error will be thrown if the task is nil
// An error will be thrown if the frequency is between 0 and MinimumFrequency (exclusive)
func NewAlarm(t time.Time, task Task, frequency time.Duration) (*Alarm, error) {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return nil, errors.New("Time must not be in the past")
	}

	if task == nil {
		return nil, errors.New("Task must not be nil.")
	}

	err := validateFrequency(frequency)
	if err != nil {
		return nil, err
	}

	return &Alarm{
		NextTime:   t,
		cancelChan: make(chan struct{}),
		task:       task,
		frequency:  frequency,
	}, nil
}

func validateFrequency(frequency time.Duration) error {
	if frequency != 0 && frequency < MinimumFrequency {
		return fmt.Errorf("Frequency must be 0 or greater than %v", MinimumFrequency)
	}
	return nil
}

// Update will change the time the alarm will finish at.
// If the alarm is currently running it will be canceled.
// An error will be thrown if time is not in the future
// An error will be thrown if the frequency is between 0 and MinimumFrequency (exclusive)
func (a *Alarm) Update(t time.Time, frequency time.Duration) error {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return errors.New("Time must not be in the past")
	}

	err := validateFrequency(frequency)
	if err != nil {
		return err
	}

	if a.running {
		a.Cancel()
	}

	a.NextTime = t
	a.cancelChan = make(chan struct{})

	return nil
}

// Cancel will cancel the alarm if it is running
// It will return an error if the alarm is not running.
func (a *Alarm) Cancel() error {
	if !a.running {
		return errors.New("Alarm not running")
	}

	close(a.cancelChan)
	return nil
}

// Start will block until either the timer goes off or the Alarm is canceled
func (a *Alarm) Start() <-chan error {
	resultChan := make(chan error)

	go func() {
		for {
			a.running = true
			durationUntilNext := a.NextTime.Sub(time.Now())
			select {
			case <-time.After(durationUntilNext):
				fmt.Printf("Alarm time has gone off\n")
				err := a.task.Run()
				if err != nil {
					fmt.Printf("Error running task: %v\n", err)
					resultChan <- err
					a.running = false
					close(resultChan)
					return
				}
				if a.frequency == 0 {
					a.running = false
					close(resultChan)
					return
				}
				a.NextTime = time.Now().Add(a.frequency)
			case <-a.cancelChan:
				fmt.Printf("Alarm canceled\n")
				a.running = false
				close(resultChan)
				return
			}
		}
	}()
	return resultChan
}
