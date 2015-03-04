package v0

import (
	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

func NewSubrouter(
	parent *mux.Router,
	taskRegistry registry.TaskRegistry,
	c *cron.Cron,
	logger lager.Logger) *mux.Router {

	r := parent.PathPrefix("/v0").Subrouter()

	tasks := r.PathPrefix("/tasks").Subrouter()
	tasks.Handle("/", tasksGetHandler(taskRegistry, logger)).Methods("GET")
	tasks.Handle("/", tasksCreateHandler(taskRegistry, logger, c)).Methods("POST")
	tasks.Handle("/{id}", taskGetHandler(taskRegistry, logger)).Methods("GET")
	tasks.Handle("/{id}", taskUpdateHandler(taskRegistry, logger, c)).Methods("PUT")
	tasks.Handle("/{id}", taskDeleteHandler(taskRegistry, logger, c)).Methods("DELETE")

	return r
}
