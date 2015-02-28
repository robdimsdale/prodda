package v0

import (
	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/registry"
)

func NewSubrouter(parent *mux.Router, prodRegistry registry.ProdRegistry) *mux.Router {
	r := parent.PathPrefix("/v0").Subrouter()

	prods := r.PathPrefix("/prods").Subrouter()
	prods.Handle("/", prodsGetHandler(prodRegistry)).Methods("GET")
	prods.Handle("/", prodsCreateHandler(prodRegistry)).Methods("POST")
	prods.Handle("/{id}", prodGetHandler(prodRegistry)).Methods("GET")

	return r
}
