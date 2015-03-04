package domain

import (
	"fmt"
	"time"

	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

const (
	MinimumTaskFrequency = time.Duration(1 * time.Minute)
)

type Task interface {
	ID() int
	SetID(id int) error

	Schedule() string
	SetSchedule(schedule string)

	EntryID() cron.EntryID
	SetEntryID(id cron.EntryID)

	Run()

	AsJSON() TaskJSON
}

type BaseTask struct {
	id       int
	schedule string
	logger   lager.Logger
	entryID  cron.EntryID
}

func (t BaseTask) ID() int {
	return t.id
}

func (t *BaseTask) SetID(id int) error {
	if t.id != 0 {
		return fmt.Errorf("Task already has an ID: %d", t.id)
	}
	t.id = id
	return nil
}

func (t BaseTask) Schedule() string {
	return t.schedule
}

func (t *BaseTask) SetSchedule(schedule string) {
	t.schedule = schedule
}

func (t BaseTask) EntryID() cron.EntryID {
	return t.entryID
}

func (t *BaseTask) SetEntryID(entryID cron.EntryID) {
	t.entryID = entryID
}

type TaskJSON interface{}

type BaseTaskJson struct {
	ID       int          `json:"id"`
	Schedule string       `json:"schedule"`
	EntryID  cron.EntryID `json:"entryID"`
	Type     string       `json:"type"`
}
