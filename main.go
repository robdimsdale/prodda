package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mfine30/prodda/api"
	"github.com/pivotal-golang/lager"
)

var (
	username string
	password string
)

func main() {
	port := os.Getenv("PORT")
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	logger := lager.NewLogger("Prodda")
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.INFO)
	logger.RegisterSink(sink)
	r := api.NewRouter(logger, username, password)

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe("0.0.0.0:"+port, r)
	}()

	fmt.Printf("Prodda started on port: %s\n", port)
	log.Fatal(<-errChan)
}
