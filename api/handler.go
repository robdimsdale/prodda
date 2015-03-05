package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prodda/prodda/api/middleware"
	"github.com/prodda/prodda/api/v0"
	"github.com/prodda/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

var HomeHandleFunc = homeHandleFunc

func NewHandler(
	logger lager.Logger,
	username, password string,
	taskRegistry registry.TaskRegistry,
	c *cron.Cron) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandleFunc)
	api := r.PathPrefix("/api").Subrouter()
	v0.NewSubrouter(api, taskRegistry, c, logger)

	return middleware.Chain{
		middleware.NewPanicRecovery(logger),
		middleware.NewLogger(logger),
		middleware.NewBasicAuth(username, password),
	}.Wrap(r)
}

func homeHandleFunc(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}
