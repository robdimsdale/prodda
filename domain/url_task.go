package domain

import (
	"io/ioutil"
	"net/http"

	"github.com/pivotal-golang/lager"
)

const (
	URLGetTaskType = "url-get"
)

type URLGetTask struct {
	BaseTask
	url string
}

type URLGetTaskJSON struct {
	BaseTaskJson
	URL string `json:"url"`
}

func NewURLGetTask(schedule, url string, logger lager.Logger) *URLGetTask {
	t := &URLGetTask{
		url: url,
	}

	t.SetSchedule(schedule)
	t.logger = logger

	return t
}

func (t URLGetTask) Run() {
	t.logger.Info("Task started", lager.Data{"task": t.AsJSON()})

	t.execute()

	t.logger.Info("Task completed", lager.Data{"task": t.AsJSON()})
	return
}

func (t URLGetTask) execute() {
	resp, err := http.Get(t.url)
	if err != nil {
		t.logger.Info(
			"Task encountered error",
			lager.Data{"task": t.AsJSON(), "err": err.Error()},
		)
		return
	}

	if resp == nil {
		t.logger.Info(
			"Task received nil response",
			lager.Data{"task": t.AsJSON()},
		)
		return
	}

	if resp.Body == nil {
		t.logger.Info(
			"Task received nil response body",
			lager.Data{"task": t.AsJSON()},
		)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.logger.Info(
			"Task encountered error reading response body",
			lager.Data{"task": t.AsJSON(), "err": err.Error()},
		)
		return
	}

	if body == nil {
		t.logger.Info(
			"Task response body nil",
			lager.Data{"task": t.AsJSON()},
		)
		return
	}

	t.logger.Info(
		"Task response",
		lager.Data{"response.body": string(body)},
	)

}

func (t URLGetTask) AsJSON() TaskJSON {
	asJson := URLGetTaskJSON{
		URL: t.url,
	}

	asJson.Type = URLGetTaskType
	asJson.ID = t.ID()
	asJson.Schedule = t.Schedule()
	asJson.EntryID = t.EntryID()

	return asJson
}
