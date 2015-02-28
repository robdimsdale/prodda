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
	"github.com/mfine30/prodda/timer"
)

type prodsCreateRequestBody struct {
	Time      time.Time `json:"time"`
	Token     string    `json:"token"`
	BuildID   uint      `json:"buildID"`
	Frequency string    `json:"frequency"`
}

func prodGetHandler(registry registry.ProdRegistry) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		idString := path.Base(r.URL.String())
		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		prod, err := registry.ByID(id)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		body, err := json.Marshal(prod.AsJSON())
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func prodsGetHandler(registry registry.ProdRegistry) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		allProds, err := registry.All()
		prodsJSON := make([]domain.ProdJSON, len(allProds))
		for i, _ := range allProds {
			prodsJSON[i] = allProds[i].AsJSON()
		}
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		body, err := json.Marshal(prodsJSON)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		fmt.Fprintf(rw, string(body))
	})
}

func prodsCreateHandler(registry registry.ProdRegistry) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var b prodsCreateRequestBody

		err := decoder.Decode(&b)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		err = validateProdRequestBody(b)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %s\n", err)
			return
		}

		// We parse the frequency separately as json decoding does not offer
		// the flexibility that ParseDuration does.
		frequency, err := time.ParseDuration(b.Frequency)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		task := domain.NewTravisTask(b.Token, b.BuildID)
		prod, err := domain.NewProd(b.Time, task, frequency)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		err = registry.Add(prod)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}

		scheduler, err := timer.NewScheduler(prod)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		fmt.Printf("Scheduler created\n")

		go scheduler.Start()

		fmt.Fprintln(rw, "New prod created")
	})
}

func validateProdRequestBody(b prodsCreateRequestBody) error {
	if b.Token == "" {
		return errors.New("Token must be provided")
	}

	if b.BuildID == 0 {
		return errors.New("BuildID must be provided")
	}
	return nil
}
