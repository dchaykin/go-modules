package apis

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dchaykin/go-modules/auth"
	log "github.com/dchaykin/go-modules/log"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			token := strings.TrimPrefix(token, "Bearer ")
			user, err := auth.GetUserIdentity(token, os.Getenv("AUTH_SECRET"))
			if err != nil {
				log.Info("Invalid user token: %v", err)
			} else {
				userData, err := json.Marshal(user)
				if err != nil {
					log.Info("Invalid user structure: %v", err)
					return
				}
				reqClone := r.Clone(r.Context())
				reqClone.Header.Set("X-User-Info", string(userData))
				next.ServeHTTP(w, reqClone)
				return
			}
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
