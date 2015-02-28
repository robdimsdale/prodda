package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/api/middleware"
	"github.com/mfine30/prodda/api/v0"
	"github.com/pivotal-golang/lager"
)

var HomeHandleFunc = homeHandleFunc

func NewHandler(logger lager.Logger, username, password string) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandleFunc)
	api := r.PathPrefix("/api").Subrouter()
	v0.NewSubrouter(api)

	return middleware.Chain{
		middleware.NewPanicRecovery(logger),
		middleware.NewLogger(logger),
		middleware.NewBasicAuth(username, password),
	}.Wrap(r)
}

func homeHandleFunc(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}
