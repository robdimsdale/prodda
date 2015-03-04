package domain

import (
	"time"

	"github.com/pivotal-golang/lager"
)

const (
	NoOpTaskType = "no-op"
)

type NoOpTask struct {
	BaseTask
	sleepDuration time.Duration
}

type NoOpTaskJSON struct {
	BaseTaskJson
	SleepDuration string `json:"sleepDuration"`
}

func NewNoOpTask(schedule string, sleepDuration time.Duration, logger lager.Logger) *NoOpTask {
	t := &NoOpTask{
		sleepDuration: sleepDuration,
	}

	t.logger = logger
	t.SetSchedule(schedule)

	return t
}

func (t NoOpTask) Run() {
	t.logger.Info("Task started", lager.Data{"task": t.AsJSON()})
	time.Sleep(t.sleepDuration)

	t.logger.Info("Task completed", lager.Data{"task": t.AsJSON()})
	return
}

func (t NoOpTask) AsJSON() TaskJSON {
	asJson := NoOpTaskJSON{
		SleepDuration: t.sleepDuration.String(),
	}

	asJson.Type = NoOpTaskType
	asJson.ID = t.ID()
	asJson.Schedule = t.Schedule()
	asJson.EntryID = t.EntryID()

	return asJson
}
