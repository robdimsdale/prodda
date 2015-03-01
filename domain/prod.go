package domain

import (
	"fmt"
	"time"
)

const (
	MiminumProdFrequency = time.Duration(1 * time.Minute)
)

type Prod struct {
	ID       int
	Task     Task
	Schedule string
}

type ProdJSON struct {
	ID       int      `json:"ID"`
	Task     TaskJSON `json:"task"`
	Schedule string   `json:"schedule"`
}

// NewProd creates a prod
func NewProd(task Task, schedule string) (*Prod, error) {
	return &Prod{
		Task:     task,
		Schedule: schedule,
	}, nil
}

func validateFrequency(frequency time.Duration) error {
	if frequency == 0 || frequency >= MiminumProdFrequency {
		return nil
	}
	return fmt.Errorf("Frequency must either be 0 or greater than %v", MiminumProdFrequency)
}

func (p Prod) Run() {
	p.Task.Run()
}

// Update will change the time the prod will finish at.
func (p *Prod) Update(schedule string) error {
	p.Schedule = schedule
	return nil
}

func (p Prod) AsJSON() ProdJSON {
	return ProdJSON{
		ID:       p.ID,
		Task:     p.Task.AsJSON(),
		Schedule: p.Schedule,
	}
}
