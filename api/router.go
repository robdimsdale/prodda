package api

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/api/v0"
)

func NewRouter(username, password string) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	api := r.PathPrefix("/api").Subrouter()
	v0.NewSubrouter(api)

	return authMiddleware(r, username, password)
}

func homeHandleFunc(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}

func authMiddleware(h http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.URL.Path == "/" {
			h.ServeHTTP(rw, req)
			return
		}

		u, p, ok := req.BasicAuth()
		if ok &&
			secureCompare(u, username) &&
			secureCompare(p, password) {
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
