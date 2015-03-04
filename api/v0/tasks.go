package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/mfine30/prodda/domain"
	"github.com/mfine30/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

type tasksCreateUpdateRequestBody struct {
	Schedule string                 `json:"schedule"`
	Task     map[string]interface{} `json:"task"`
}

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
		var b tasksCreateUpdateRequestBody

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
		decoder := json.NewDecoder(r.Body)
		var b tasksCreateUpdateRequestBody

		err := decoder.Decode(&b)
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

		taskTypeRaw := b.Task["type"]
		if taskTypeRaw == nil {
			err := errors.New("Task type must be provided")
			logger.Info("Failed to create task", lager.Data{"err": err.Error()})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		taskType := taskTypeRaw.(string)

		var task domain.Task
		switch taskType {
		case domain.TravisTaskType:
			task, err = createTravisTaskConfig(b, logger)
			if err != nil {
				logger.Info("Failed to create Travis task", lager.Data{"err": err.Error()})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		case domain.NoOpTaskType:
			task, err = createNoOpTaskConfig(b, logger)
			if err != nil {
				logger.Info("Failed to create NoOp task", lager.Data{"err": err.Error()})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		case domain.URLGetTaskType:
			task, err = createURLGetTaskConfig(b, logger)
			if err != nil {
				logger.Info("Failed to create URLGet task", lager.Data{"err": err.Error()})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		default:
			err := fmt.Errorf("Unrecognized task type: %s", taskType)
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

		body, err := json.Marshal(task.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize task", err, lager.Data{"task": task.AsJSON()})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		logger.Info("Task created", lager.Data{"task": task.AsJSON()})

		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, string(body))

	})
}

func createTravisTaskConfig(b tasksCreateUpdateRequestBody, logger lager.Logger) (*domain.TravisTask, error) {
	task := b.Task
	tokenRaw := task["token"]
	buildIDRaw := task["buildID"]

	if tokenRaw == nil {
		return nil, errors.New("Token must be provided")
	}

	if buildIDRaw == nil {
		return nil, errors.New("BuildID must be provided")
	}

	var token string
	switch tokenRaw.(type) {
	case string:
		token = tokenRaw.(string)
	default:
		return nil, fmt.Errorf("Cannot parse token: %v", tokenRaw)
	}

	var buildID64 float64
	switch buildIDRaw.(type) {
	case float64:
		buildID64 = buildIDRaw.(float64)
	default:
		return nil, fmt.Errorf("Cannot parse buildID: %v", buildIDRaw)
	}

	buildID := uint(buildID64)

	return domain.NewTravisTask(b.Schedule, token, buildID, logger), nil
}

func createNoOpTaskConfig(b tasksCreateUpdateRequestBody, logger lager.Logger) (*domain.NoOpTask, error) {
	task := b.Task
	sleepDurationRaw := task["sleepDuration"]

	if sleepDurationRaw == nil {
		return nil, errors.New("Sleep duration must be provided")
	}

	var sleepDurationString string
	switch sleepDurationRaw.(type) {
	case string:
		sleepDurationString = sleepDurationRaw.(string)
	default:
		return nil, fmt.Errorf("Cannot parse sleep duration: %v", sleepDurationRaw)
	}

	sleepDuration, err := time.ParseDuration(sleepDurationString)
	if err != nil {
		return nil, err
	}

	return domain.NewNoOpTask(b.Schedule, sleepDuration, logger), nil
}

func createURLGetTaskConfig(b tasksCreateUpdateRequestBody, logger lager.Logger) (*domain.URLGetTask, error) {
	task := b.Task
	urlRaw := task["url"]

	if urlRaw == nil {
		return nil, errors.New("URL must be provided")
	}

	var urlString string
	switch urlRaw.(type) {
	case string:
		urlString = urlRaw.(string)
	default:
		return nil, fmt.Errorf("Cannot parse sleep duration: %v", urlRaw)
	}

	return domain.NewURLGetTask(b.Schedule, urlString, logger), nil
}
