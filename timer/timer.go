package timer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mfine30/prodda/domain"
)

type Scheduler struct {
	mutex      *sync.Mutex
	running    bool
	cancelChan chan struct{}
	prod       *domain.Prod
}

// NewScheduler creates a scheduler
func NewScheduler(p *domain.Prod) (*Scheduler, error) {
	return &Scheduler{
		prod:       p,
		cancelChan: make(chan struct{}),
		mutex:      &sync.Mutex{},
	}, nil
}

func (a Scheduler) NextTime() time.Time {
	return a.prod.NextTime
}

// Update will change the time the at which prod will be scheduled.
// If the scheduler is currently running it will be canceled,
// and Start will need to be called again.
func (a *Scheduler) Update(t time.Time, frequency time.Duration) error {
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

// Cancel will cancel the scheduler if it is running
// It will return an error if the scheduler is not running.
func (a *Scheduler) Cancel() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if !a.running {
		return errors.New("Scheduler not running")
	}

	close(a.cancelChan)
	return nil
}

// Start will block until either the timer goes off or the scheduler is canceled
func (a *Scheduler) Start() <-chan error {
	resultChan := make(chan error)

	go func() {
		for {
			a.mutex.Lock()
			a.running = true
			a.mutex.Unlock()

			durationUntilNext := a.prod.NextTime.Sub(time.Now())
			select {
			case <-time.After(durationUntilNext):
				fmt.Printf("Scheduler time has gone off\n")
				err := a.prod.Run()
				if err != nil {
					fmt.Printf("Error running prod: %v\n", err)

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
				fmt.Printf("Scheduler canceled\n")

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
