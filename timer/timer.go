package timer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mfine30/prodda/domain"
)

type Alarm struct {
	mutex      *sync.Mutex
	running    bool
	cancelChan chan struct{}
	prod       *domain.Prod
}

// NewAlarm creates an alarm
func NewAlarm(p *domain.Prod) (*Alarm, error) {
	return &Alarm{
		prod:       p,
		cancelChan: make(chan struct{}),
		mutex:      &sync.Mutex{},
	}, nil
}

func (a Alarm) NextTime() time.Time {
	return a.prod.NextTime
}

// Update will change the time the alarm will finish at.
// If the alarm is currently running it will be canceled.
func (a *Alarm) Update(t time.Time, frequency time.Duration) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.prod.Update(t, frequency)
	if err != nil {
		return err
	}

	if a.running {
		a.Cancel()
	}

	a.cancelChan = make(chan struct{})

	return nil
}

// Cancel will cancel the alarm if it is running
// It will return an error if the alarm is not running.
func (a *Alarm) Cancel() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

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
			a.mutex.Lock()
			a.running = true
			a.mutex.Unlock()

			durationUntilNext := a.prod.NextTime.Sub(time.Now())
			select {
			case <-time.After(durationUntilNext):
				fmt.Printf("Alarm time has gone off\n")
				err := a.prod.Run()
				if err != nil {
					fmt.Printf("Error running task: %v\n", err)

					a.mutex.Lock()
					a.running = false
					a.mutex.Unlock()

					resultChan <- err
					close(resultChan)
					return
				}
				if a.prod.Frequency == 0 {
					a.mutex.Lock()
					a.running = false
					a.mutex.Unlock()

					close(resultChan)
					return
				}
				a.prod.NextTime = time.Now().Add(a.prod.Frequency)
			case <-a.cancelChan:
				fmt.Printf("Alarm canceled\n")

				a.mutex.Lock()
				a.running = false
				a.mutex.Unlock()

				close(resultChan)
				return
			}
		}
	}()
	return resultChan
}
