package main

import (
	"os"
	"strconv"

	"github.com/mfine30/prodda/api"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
)

var (
	username string
	password string
)

func main() {
	logger := lager.NewLogger("Prodda")
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.INFO)
	logger.RegisterSink(sink)

	portEnv := os.Getenv("PORT")
	port64, err := strconv.ParseUint(portEnv, 10, 0)
	port := uint(port64)
	if err != nil {
		logger.Fatal("Cannot parse port from environment", err, lager.Data{"PORT": portEnv})
	}
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	handler := api.NewHandler(logger, username, password)

	group := grouper.NewParallel(os.Kill, grouper.Members{
		grouper.Member{"api", api.NewRunner(port, handler, logger)},
	})
	process := ifrit.Invoke(group)

	logger.Info("Prodda started")
	err = <-process.Wait()
	if err != nil {
		logger.Fatal("Error running prodda", err)
	}
}
