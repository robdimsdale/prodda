package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/mfine30/prodda/domain"
	"github.com/mfine30/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

func taskGetHandler(registry registry.TaskRegistry, logger lager.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Info("Failed to get task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		task, err := registry.ByID(id)
		if err != nil {
			logger.Error("Failed to find existing task in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if task == nil {
			logger.Info("Task not found in registry", lager.Data{"ID": id})
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: task not found for ID: %d\n", id)
			return
		}

		body, err := json.Marshal(task.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize task", err, lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func taskUpdateHandler(registry registry.TaskRegistry, logger lager.Logger, c *cron.Cron) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Info("Failed to update task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		task, err := registry.ByID(id)
		if err != nil {
			logger.Error("Failed to find existing task in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if task == nil {
			logger.Info("Task not found in registry", lager.Data{"ID": id})
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: task not found for ID: %d\n", id)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var b domain.BaseTaskJson

		err = decoder.Decode(&b)
		if err != nil {
			logger.Error("Failed to update task", err, lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		task.SetSchedule(b.Schedule)
		c.Remove(task.EntryID())
		entryID, err := c.AddJob(task.Schedule(), task)
		if err != nil {
			logger.Error(
				"Failed to schedule task",
				err,
				lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		task.SetEntryID(entryID)

		task, err = registry.Update(task)
		if err != nil {
			logger.Error("Failed to update task in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		logger.Info("task updated", lager.Data{"task": task.AsJSON()})

		body, err := json.Marshal(task.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize task", err, lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func taskDeleteHandler(registry registry.TaskRegistry, logger lager.Logger, c *cron.Cron) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Info("Failed to delete task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		task, err := registry.ByID(id)
		if err != nil {
			logger.Error("Failed to find existing task in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if task == nil {
			logger.Info("Task not found in registry", lager.Data{"ID": id})
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: task not found for ID: %d\n", id)
			return
		}

		c.Remove(task.EntryID())

		err = registry.Remove(task)
		if err != nil {
			logger.Error("Failed to remove task from registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		logger.Info("task deleted", lager.Data{"task": task.AsJSON()})
	})
}

func tasksGetHandler(registry registry.TaskRegistry, logger lager.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		allTasks, err := registry.All()
		if err != nil {
			logger.Error("Failed to get tasks from registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		tasksJSON := make([]domain.TaskJSON, len(allTasks))
		for i, _ := range allTasks {
			tasksJSON[i] = allTasks[i].AsJSON()
		}

		body, err := json.Marshal(tasksJSON)
		if err != nil {
			logger.Error("Failed to serialize tasks", err, lager.Data{"tasks": tasksJSON})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func tasksCreateHandler(registry registry.TaskRegistry, logger lager.Logger, c *cron.Cron) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Info("Failed to create task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		var b domain.BaseTaskJson

		err = json.Unmarshal(body, &b)
		if err != nil {
			logger.Info("Failed to create task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if b.Schedule == "" {
			err := errors.New("Schedule must be provided")
			logger.Info("Failed to create task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if b.Type == "" {
			err := errors.New("Task type must be provided")
			logger.Info("Failed to create task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		var task domain.Task
		switch b.Type {
		case domain.TravisTaskType:
			task, err = createTravisTaskConfig(body, logger)
			if err != nil {
				logger.Info("Failed to create Travis task", lager.Data{"err": err.Error()})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		case domain.NoOpTaskType:
			task, err = createNoOpTaskConfig(body, logger)
			if err != nil {
				logger.Info("Failed to create NoOp task", lager.Data{"err": err.Error()})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		case domain.URLGetTaskType:
			task, err = createURLGetTaskConfig(body, logger)
			if err != nil {
				logger.Info("Failed to create URLGet task", lager.Data{"err": err.Error()})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		default:
			err := fmt.Errorf("Unrecognized task type: %s", b.Type)
			logger.Info("Failed to create task", lager.Data{"err": err.Error()})
			rw.WriteHeader(httpUnprocessableEntity)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		entryID, err := c.AddJob(b.Schedule, task)
		if err != nil {
			logger.Error(
				"Failed to schedule task",
				err,
				lager.Data{"schedule": b.Schedule, "task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		task.SetEntryID(entryID)

		err = registry.Add(task)
		if err != nil {
			logger.Error("Failed to add task to registry", err, lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		responseBody, err := json.Marshal(task.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize task", err, lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		logger.Info("Task created", lager.Data{"task": task.AsJSON()})

		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, string(responseBody))

	})
}

func createTravisTaskConfig(b []byte, logger lager.Logger) (*domain.TravisTask, error) {
	var task travisTaskConfig
	err := json.Unmarshal(b, &task)
	if err != nil {
		return nil, err
	}

	if task.Token == "" {
		return nil, errors.New("Token must be provided")
	}

	if task.BuildID == 0 {
		return nil, errors.New("BuildID must be provided")
	}

	return domain.NewTravisTask(task.Schedule, task.Token, task.BuildID, logger), nil
}

type travisTaskConfig struct {
	domain.TravisTaskJSON
	Token string `json:"token"`
}

func createNoOpTaskConfig(b []byte, logger lager.Logger) (*domain.NoOpTask, error) {
	var task domain.NoOpTaskJSON
	err := json.Unmarshal(b, &task)
	if err != nil {
		return nil, err
	}

	var sleepDuration time.Duration
	if task.SleepDuration != "" {
		sleepDuration, err = time.ParseDuration(task.SleepDuration)
		if err != nil {
			return nil, err
		}
	} else {
		sleepDuration = 0
	}

	return domain.NewNoOpTask(task.Schedule, sleepDuration, logger), nil
}

func createURLGetTaskConfig(b []byte, logger lager.Logger) (*domain.URLGetTask, error) {
	var task domain.URLGetTaskJSON
	err := json.Unmarshal(b, &task)
	if err != nil {
		return nil, err
	}

	if task.URL == "" {
		return nil, errors.New("URL must be provided")
	}

	return domain.NewURLGetTask(task.Schedule, task.URL, logger), nil
}
