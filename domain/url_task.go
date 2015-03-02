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
	url    string
	logger lager.Logger
}

type URLGetTaskJSON struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func NewURLGetTask(url string, logger lager.Logger) *URLGetTask {
	return &URLGetTask{
		url:    url,
		logger: logger,
	}
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
	return URLGetTaskJSON{
		Type: URLGetTaskType,
		URL:  t.url,
	}
}
