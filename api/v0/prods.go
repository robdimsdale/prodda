package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/mfine30/prodda/domain"
	"github.com/mfine30/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

type prodsCreateUpdateRequestBody struct {
	Token    string `json:"token"`
	BuildID  uint   `json:"buildID"`
	Schedule string `json:"schedule"`
}

func prodGetHandler(registry registry.ProdRegistry, logger lager.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
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
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var b prodsCreateUpdateRequestBody

		err = decoder.Decode(&b)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod.Schedule = b.Schedule
		c.Remove(cron.EntryID(prod.ID))
		c.AddJob(prod.Schedule, prod.Task)

		prod, err = registry.Update(prod)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		logger.Info("prod updated", lager.Data{"prod": prod})

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
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
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		c.Remove(cron.EntryID(prod.ID))

		err = registry.Remove(prod)
		if err != nil {
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
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		err = validateProdRequestBody(b)
		if err != nil {
			rw.WriteHeader(httpUnprocessableEntity)
			fmt.Fprintf(rw, "ERROR: %s\n", err)
			return
		}

		task := domain.NewTravisTask(b.Token, b.BuildID, logger)
		prod, err := domain.NewProd(task, b.Schedule)
		if err != nil {
			rw.WriteHeader(httpUnprocessableEntity)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		err = registry.Add(prod)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		c.AddJob(prod.Schedule, prod.Task)

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		logger.Info("prod created", lager.Data{"prod": prod})

		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, string(body))

	})
}

func validateProdRequestBody(b prodsCreateUpdateRequestBody) error {
	if b.Token == "" {
		return errors.New("Token must be provided")
	}

	if b.BuildID == 0 {
		return errors.New("BuildID must be provided")
	}
	return nil
}
