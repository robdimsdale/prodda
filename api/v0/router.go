package v0

import "github.com/gorilla/mux"

func NewSubrouter(parent *mux.Router) *mux.Router {
	r := parent.PathPrefix("/v0").Subrouter()
	prods := r.PathPrefix("/prods").Subrouter()
	prods.Methods("POST").Handler(prodsCreateHandler())
	return r
}
