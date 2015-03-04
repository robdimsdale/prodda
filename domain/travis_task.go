package domain

import (
	"github.com/mfine30/prodda/client"
	"github.com/pivotal-golang/lager"
)

const (
	TravisTaskType = "travis-re-run"
)

type TravisTask struct {
	BaseTask
	client  *client.Travis
	token   string
	buildID uint
}

type TravisTaskJSON struct {
	BaseTaskJson
	BuildID uint `json:"buildID"`
}

func NewTravisTask(schedule, token string, buildID uint, logger lager.Logger) *TravisTask {
	t := &TravisTask{
		client:  client.NewTravisClient("https://api.travis-ci.org"),
		token:   token,
		buildID: buildID,
	}

	t.logger = logger
	t.SetSchedule(schedule)

	return t
}

func (t TravisTask) Run() {
	t.logger.Info("Task started", lager.Data{"task": t.AsJSON()})

	response, _ := t.client.TriggerBuild(t.token, t.buildID)
	t.logger.Info("Task completed", lager.Data{"task": t.AsJSON(), "response": response})
	return
}

func (t TravisTask) AsJSON() TaskJSON {
	asJson := TravisTaskJSON{
		BuildID: t.buildID,
	}

	asJson.Type = TravisTaskType
	asJson.ID = t.ID()
	asJson.Schedule = t.Schedule()
	asJson.EntryID = t.EntryID()

	return asJson
}
