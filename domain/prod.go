package domain

import (
	"fmt"
	"time"

	"gopkg.in/robfig/cron.v2"
)

const (
	MiminumProdFrequency = time.Duration(1 * time.Minute)
)

type Prod struct {
	ID       int
	Task     Task
	Schedule string
	EntryID  cron.EntryID
}

type ProdJSON struct {
	ID       int      `json:"ID"`
	Task     TaskJSON `json:"task"`
	Schedule string   `json:"schedule"`
}

// NewProd creates a prod
func NewProd(task Task, schedule string, entryID cron.EntryID) (*Prod, error) {
	return &Prod{
		Task:     task,
		Schedule: schedule,
		EntryID:  entryID,
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
