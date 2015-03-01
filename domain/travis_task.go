package domain

import (
	"github.com/mfine30/prodda/client"
	"github.com/pivotal-golang/lager"
)

type TravisTask struct {
	client  *client.Travis
	token   string
	buildID uint
	logger  lager.Logger
}

type TravisTaskJSON struct {
	Type    string `json:"type"`
	BuildID uint   `json:"buildID"`
}

func NewTravisTask(token string, buildID uint, logger lager.Logger) *TravisTask {
	return &TravisTask{
		client:  client.NewTravisClient("https://api.travis-ci.org"),
		token:   token,
		buildID: buildID,
		logger:  logger,
	}
}

func (t TravisTask) Run() {
	t.logger.Info("Task started", lager.Data{"task": t.AsJSON()})

	response, _ := t.client.TriggerBuild(t.token, t.buildID)
	t.logger.Info("Task completed", lager.Data{"task": t.AsJSON(), "response": response})
	return
}

func (t TravisTask) AsJSON() TaskJSON {
	return TravisTaskJSON{
		Type:    "Travis:",
		BuildID: t.buildID,
	}
}
