package domain

import (
	"errors"
	"fmt"
	"time"
)

const (
	MiminumProdFrequency = time.Duration(1 * time.Minute)
)

type Prod struct {
	Task      Task
	NextTime  time.Time
	Frequency time.Duration
}

type ProdJSON struct {
	Task      TaskJSON      `json:"task"`
	NextTime  time.Time     `json:"nextTime"`
	Frequency time.Duration `json:"duration"`
}

// NewProd creates a prod
// An error will be thrown if time is not in the future
// An error will be thrown if the task is nil
// An error will be thrown if the frequency is between 0 and MinimumFrequency (exclusive)
func NewProd(t time.Time, task Task, frequency time.Duration) (*Prod, error) {
	currentTime := time.Now()

	if t.IsZero() {
		t = currentTime
	} else if t.Before(currentTime) {
		return nil, errors.New("Time must not be in the past")
	}

	if task == nil {
		return nil, errors.New("Task must not be nil.")
	}

	err := validateFrequency(frequency)
	if err != nil {
		return nil, err
	}

	return &Prod{
		NextTime:  t,
		Task:      task,
		Frequency: frequency,
	}, nil
}

func validateFrequency(frequency time.Duration) error {
	if frequency == 0 || frequency >= MiminumProdFrequency {
		return nil
	}
	return fmt.Errorf("Frequency must either be 0 or greater than %v", MiminumProdFrequency)
}

func (p Prod) Run() error {
	return p.Task.Run()
}

// Update will change the time the prod will finish at.
// An error will be thrown if time is not in the future
// An error will be thrown if the frequency is between 0 and MinimumFrequency (exclusive)
func (p *Prod) Update(t time.Time, frequency time.Duration) error {
	currentTime := time.Now()
	if t.Before(currentTime) {
		return errors.New("Time must not be in the past")
	}

	err := validateFrequency(frequency)
	if err != nil {
		return err
	}

	p.NextTime = t
	p.Frequency = frequency
	return nil
}

func (p Prod) AsJSON() ProdJSON {
	return ProdJSON{
		Task:      p.Task.AsJSON(),
		NextTime:  p.NextTime,
		Frequency: p.Frequency,
	}
}
