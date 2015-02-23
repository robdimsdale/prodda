package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/timer"
)

func main() {
	port := os.Getenv("PORT")

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler)
	prods := r.Path("/prods").Subrouter()
	prods.Methods("POST").HandlerFunc(prodsCreateHandler)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}

type prodsCreateRequestBody struct {
	Time    time.Time `json:"time"`
	Token   string    `json:"token"`
	BuildID uint      `json:"build_id"`
}

func prodsCreateHandler(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var b prodsCreateRequestBody

	err := decoder.Decode(&b)
	if err != nil {
		fmt.Fprintf(rw, "ERROR: %v\n", err)
		return
	}

	if b.Time.IsZero() {
		fmt.Fprintf(rw, "ERROR: Time must be provided\n")
		return
	}

	if b.Token == "" {
		fmt.Fprintf(rw, "ERROR: Token must be provided\n")
		return
	}

	if b.BuildID == 0 {
		fmt.Fprintf(rw, "ERROR: BuildID must be provided\n")
		return
	}

	task := timer.NewTravisTask(b.Token, b.BuildID)
	alarm, err := timer.NewAlarm(b.Time, task)
	if err != nil {
		fmt.Fprintf(rw, "ERROR: %v\n", err)
		return
	}
	fmt.Printf("Alarm created\n")

	go alarm.Start()

	fmt.Fprintln(rw, "New prod created")
}
