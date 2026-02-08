package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("")
		log.Printf("%s %s %s", timestamp, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
