package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/timer"
)

var (
	username string
	password string
)

func main() {
	port := os.Getenv("PORT")
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	prods := r.PathPrefix("/prods").Subrouter()
	prods.Methods("POST").Handler(prodsCreateHandler())

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe("0.0.0.0:"+port, authMiddleware(r))
	}()

	fmt.Printf("Prodda started on port: %s\n", port)
	log.Fatal(<-errChan)
}

func homeHandleFunc(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		fmt.Println("in auth middleware")
		if req.URL.Path == "/" {
			h.ServeHTTP(rw, req)
			return
		}

		u, p, ok := req.BasicAuth()
		if ok &&
			secureCompare(u, username) &&
			secureCompare(p, password) {
			fmt.Println("auth good")
			fmt.Printf("h: %+v", h)
			h.ServeHTTP(rw, req)
		} else {
			rw.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
		}
	})
}

func secureCompare(a, b string) bool {
	x := []byte(a)
	y := []byte(b)
	return subtle.ConstantTimeCompare(x, y) == 1
}

type prodsCreateRequestBody struct {
	Time      time.Time `json:"time"`
	Token     string    `json:"token"`
	BuildID   uint      `json:"build_id"`
	Frequency string    `json:"frequency"`
}

func prodsCreateHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("in prod create handler")
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

		task := timer.NewTravisTask(b.Token, b.BuildID)
		alarm, err := timer.NewAlarm(b.Time, task, frequency)
		if err != nil {
			fmt.Fprintf(rw, "ERROR: %v\n", err)
			return
		}
		fmt.Printf("Alarm created\n")

		go alarm.Start()

		fmt.Fprintln(rw, "New prod created")
	})
}

func validateProdRequestBody(b prodsCreateRequestBody) error {
	if b.Time.IsZero() {
		return errors.New("Time must be provided")
	}

	if b.Token == "" {
		return errors.New("Token must be provided")
	}

	if b.BuildID == 0 {
		return errors.New("BuildID must be provided")
	}
	return nil
}
