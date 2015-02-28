package v0

import (
	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/registry"
)

func NewSubrouter(parent *mux.Router, prodRegistry registry.ProdRegistry) *mux.Router {
	r := parent.PathPrefix("/v0").Subrouter()
	prods := r.PathPrefix("/prods").Subrouter()
	prods.Methods("GET").Handler(prodsGetHandler(prodRegistry))
	prods.Methods("POST").Handler(prodsCreateHandler(prodRegistry))
	return r
}
