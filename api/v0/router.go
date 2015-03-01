package v0

import (
	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

func NewSubrouter(
	parent *mux.Router,
	prodRegistry registry.ProdRegistry,
	c *cron.Cron,
	logger lager.Logger) *mux.Router {

	r := parent.PathPrefix("/v0").Subrouter()

	prods := r.PathPrefix("/prods").Subrouter()
	prods.Handle("/", prodsGetHandler(prodRegistry, logger)).Methods("GET")
	prods.Handle("/", prodsCreateHandler(prodRegistry, logger, c)).Methods("POST")
	prods.Handle("/{id}", prodGetHandler(prodRegistry, logger)).Methods("GET")
	prods.Handle("/{id}", prodUpdateHandler(prodRegistry, logger, c)).Methods("POST")
	prods.Handle("/{id}", prodDeleteHandler(prodRegistry, logger, c)).Methods("DELETE")

	return r
}
