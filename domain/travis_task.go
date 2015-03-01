package domain

import (
	"fmt"
	"log"

	"github.com/mfine30/prodda/client"
)

type TravisTask struct {
	client  *client.Travis
	token   string
	buildID uint
}

type TravisTaskJSON struct {
	Type    string `json:"type"`
	BuildID uint   `json:"buildID"`
}

func NewTravisTask(token string, buildID uint) *TravisTask {
	return &TravisTask{
		client:  client.NewTravisClient("https://api.travis-ci.org"),
		token:   token,
		buildID: buildID,
	}
}

func (t TravisTask) Run() {
	fmt.Printf("Travis task running\n")

	resp, _ := t.client.TriggerBuild(t.token, t.buildID)
	log.Printf("response: %+v\n", resp)
	return
}

func (t TravisTask) AsJSON() TaskJSON {
	return TravisTaskJSON{
		Type:    "Travis:",
		BuildID: t.buildID,
	}
}
