package apis

import (
	"net/http"
	"strings"

	"github.com/dchaykin/go-modules/log"
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

func AddAdminEndpoints(router *mux.Router) {
	router.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		userIdentity, err := getUserIdentityFromToken(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: " + err.Error()))
			return
		}
		if !userIdentity.IsDeveloper() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized: developer abilities are needed"))
			return
		}
		logLevel := r.URL.Query().Get("level")
		switch strings.ToLower(logLevel) {
		case "debug":
			log.SetLevel(log.LevelDebug)
		case "warn", "warning":
			log.SetLevel(log.LevelWarn)
		case "info":
			log.SetLevel(log.LevelInfo)
		case "error":
			log.SetLevel(log.LevelError)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown log level" + logLevel))
			return
		}
		w.Write([]byte("Log level is changed to " + logLevel))
	})
}
