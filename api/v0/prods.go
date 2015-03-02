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

type prodsCreateUpdateRequestBody struct {
	Schedule string                 `json:"schedule"`
	Task     map[string]interface{} `json:"task"`
}

func prodGetHandler(registry registry.ProdRegistry, logger lager.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Info("Failed to get prod", lager.Data{"err": err})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			logger.Error("Failed to find existing prod in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if prod == nil {
			logger.Info("Prod not found in registry", lager.Data{"ID": id})
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: prod not found for ID: %d\n", id)
			return
		}

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize prod", err, lager.Data{"prod": prod})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func prodUpdateHandler(registry registry.ProdRegistry, logger lager.Logger, c *cron.Cron) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Info("Failed to update prod", lager.Data{"err": err})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			logger.Error("Failed to find existing prod in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if prod == nil {
			logger.Info("Prod not found in registry", lager.Data{"ID": id})
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: prod not found for ID: %d\n", id)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var b prodsCreateUpdateRequestBody

		err = decoder.Decode(&b)
		if err != nil {
			logger.Error("Failed to update prod", err, lager.Data{"prod": prod})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod.Schedule = b.Schedule
		c.Remove(cron.EntryID(prod.ID))
		c.AddJob(prod.Schedule, prod.Task)

		prod, err = registry.Update(prod)
		if err != nil {
			logger.Error("Failed to update prod in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		logger.Info("prod updated", lager.Data{"prod": prod})

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize prod", err, lager.Data{"prod": prod})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func prodDeleteHandler(registry registry.ProdRegistry, logger lager.Logger, c *cron.Cron) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Info("Failed to delete prod", lager.Data{"err": err})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			logger.Error("Failed to find existing prod in registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if prod == nil {
			logger.Info("Prod not found in registry", lager.Data{"ID": id})
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: prod not found for ID: %d\n", id)
			return
		}

		c.Remove(cron.EntryID(prod.ID))

		err = registry.Remove(prod)
		if err != nil {
			logger.Error("Failed to remove prod from registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		logger.Info("prod deleted", lager.Data{"prod": prod})
	})
}

func prodsGetHandler(registry registry.ProdRegistry, logger lager.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		allProds, err := registry.All()
		if err != nil {
			logger.Error("Failed to get prods from registry", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prodsJSON := make([]domain.ProdJSON, len(allProds))
		for i, _ := range allProds {
			prodsJSON[i] = allProds[i].AsJSON()
		}

		body, err := json.Marshal(prodsJSON)
		if err != nil {
			logger.Error("Failed to serialize prods", err, lager.Data{"prod": prodsJSON})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func prodsCreateHandler(registry registry.ProdRegistry, logger lager.Logger, c *cron.Cron) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var b prodsCreateUpdateRequestBody

		err := decoder.Decode(&b)
		if err != nil {
			logger.Info("Failed to create task", lager.Data{"err": err})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		if b.Schedule == "" {
			err := errors.New("Schedule must be provided")
			logger.Info("Failed to create task", lager.Data{"err": err})
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		taskTypeRaw := b.Task["type"]
		if taskTypeRaw == nil {
			err := errors.New("Task type must be provided")
			logger.Info("Failed to create task", lager.Data{"err": err})
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
				logger.Info("Failed to create Travis task", lager.Data{"err": err})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		case domain.NoOpTaskType:
			task, err = createNoOpTaskConfig(b, logger)
			if err != nil {
				logger.Info("Failed to create NoOp task", lager.Data{"err": err})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		case domain.URLGetTaskType:
			task, err = createURLGetTaskConfig(b, logger)
			if err != nil {
				logger.Info("Failed to create URLGet task", lager.Data{"err": err})
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "ERROR: %v\n", err)
				return
			}
		default:
			err := fmt.Errorf("Unrecognized task type: %s", taskType)
			logger.Info("Failed to create task", lager.Data{"err": err})
			rw.WriteHeader(httpUnprocessableEntity)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := domain.NewProd(task, b.Schedule)
		if err != nil {
			logger.Info("Failed to create Prod", lager.Data{"err": err})
			rw.WriteHeader(httpUnprocessableEntity)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		err = registry.Add(prod)
		if err != nil {
			logger.Error("Failed to add prod to registry", err, lager.Data{"prod": prod})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		c.AddJob(prod.Schedule, prod.Task)

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
			logger.Error("Failed to serialize prod", err, lager.Data{"prod": prod})
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		logger.Info("prod created", lager.Data{"prod": prod})

		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, string(body))

	})
}

func createTravisTaskConfig(b prodsCreateUpdateRequestBody, logger lager.Logger) (*domain.TravisTask, error) {
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

	return domain.NewTravisTask(token, buildID, logger), nil
}

func createNoOpTaskConfig(b prodsCreateUpdateRequestBody, logger lager.Logger) (*domain.NoOpTask, error) {
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

	return domain.NewNoOpTask(sleepDuration, logger), nil
}

func createURLGetTaskConfig(b prodsCreateUpdateRequestBody, logger lager.Logger) (*domain.URLGetTask, error) {
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

	return domain.NewURLGetTask(urlString, logger), nil
}
