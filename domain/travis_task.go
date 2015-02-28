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

func NewTravisTask(token string, buildID uint) *TravisTask {
	return &TravisTask{
		client:  client.NewTravisClient("https://api.travis-ci.org"),
		token:   token,
		buildID: buildID,
	}
}

func (t TravisTask) Run() error {
	fmt.Printf("Travis task running\n")

	resp, err := t.client.TriggerBuild(t.token, t.buildID)
	if err != nil {
		return err
	}
	log.Printf("response: %+v\n", resp)
	return nil
}
