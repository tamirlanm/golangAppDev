package handler

import (
	"log"
	"net/http"
	"time"
)

const validAPIKey = "my-secret-key"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] methods=%s endpoints=%s duration=%s", start.Format(time.RFC3339), r.Method, r.RequestURI, time.Since(start))
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key != validAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized: invalid or missing API key"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
