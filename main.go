package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/timer"
)

var (
	travisToken string
)

func main() {
	port := os.Getenv("PORT")
	travisToken = os.Getenv("TRAVIS_TOKEN")

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler)
	prods := r.Path("/prods").Subrouter()
	prods.Methods("POST").HandlerFunc(prodsCreateHandler)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}

func prodsCreateHandler(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Printf("Query: %+v\n", query)
	timeQuery := query["time"]
	if timeQuery == nil {
		log.Println("Time not present in URL params")
		fmt.Fprintln(rw, "ERROR: time must be present in URL params")
		return
	}

	fmt.Printf("time-now: %v\n", time.Now().Format(time.RFC3339))
	fmt.Printf("time [0]: %s\n", timeQuery[0])
	t, err := time.Parse(time.RFC3339, timeQuery[0])
	if err != nil {
		log.Printf("Cannot parse time from %+v\n", timeQuery[0])
		fmt.Fprintf(rw, "ERROR: Cannot parse time from %+s\n", timeQuery[0])
		return
	}

	task := timer.NewTravisTask(travisToken)
	alarm, err := timer.NewAlarm(t, task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Alarm created\n")

	go alarm.Start()

	fmt.Fprintln(rw, "New prod created")
}
