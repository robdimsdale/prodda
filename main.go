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
	dingAt := time.Now().Add(5 * time.Second)
	task := timer.NewTravisTask(travisToken)
	alarm, err := timer.NewAlarm(dingAt, task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Alarm created\n")

	go alarm.Start()

	fmt.Fprintln(rw, "New prod created")
}
