package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Build struct {
	Id int `json:"id"`
}

type Travis struct {
	url string
}

func NewTravisClient(apiServer string) Travis {
	return Travis{apiServer}
}

func (t *Travis) GetBuilds(user, repo string) (*[]Build, error) {
	resp, err := http.Get(fmt.Sprintf("%s/repos/%s/%s/builds", t.url, user, repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var builds *[]Build = &[]Build{}
	err = json.Unmarshal(respBody, builds)
	if err != nil {
		return nil, err
	}
	return builds, err
}

type RestartNotice struct {
	Notice string `json:"notice"`
}

type RestartResponse struct {
	Result bool            `json:"result"`
	Flash  []RestartNotice `json:"flash"`
}

func (t *Travis) TriggerBuild(user, repo, travisToken string, buildId int) (string, error) {
	URL := fmt.Sprintf("%s/requests", t.url)
	formBody := fmt.Sprintf(`{"build_id": %d}`, buildId)
	body := ioutil.NopCloser(strings.NewReader(formBody))

	request, err := http.NewRequest(
		"POST",
		URL,
		body)

	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", fmt.Sprintf("token %s", travisToken))
	request.Header.Set("Accept", "application/json; version=2")
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var restartResponse RestartResponse
	err = json.Unmarshal(respBody, &restartResponse)
	if err != nil {
		return "", err
	}

	return restartResponse.Flash[0].Notice, nil
}
