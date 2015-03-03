package domain

import (
	"time"

	"github.com/pivotal-golang/lager"
)

const (
	NoOpTaskType = "no-op"
)

type NoOpTask struct {
	sleepDuration time.Duration
	logger        lager.Logger
}

type NoOpTaskJSON struct {
	Type          string        `json:"type"`
	SleepDuration time.Duration `json:"sleepDuration"`
}

func NewNoOpTask(sleepDuration time.Duration, logger lager.Logger) *NoOpTask {
	return &NoOpTask{
		sleepDuration: sleepDuration,
		logger:        logger,
	}
}

func (t NoOpTask) Run() {
	t.logger.Info("Task started", lager.Data{"task": t.AsJSON()})
	time.Sleep(t.sleepDuration)

	t.logger.Info("Task completed", lager.Data{"task": t.AsJSON()})
	return
}

func (t NoOpTask) AsJSON() TaskJSON {
	return NoOpTaskJSON{
		Type:          NoOpTaskType,
		SleepDuration: t.sleepDuration,
	}
}
