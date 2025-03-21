package apis

import (
	"net/http"

	"github.com/gorilla/mux"
)

func AddStandardEndpoints(router *mux.Router) {
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	router.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		// TODO: make more here
		w.Write([]byte("Ready"))
	})
}
