package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// #nosec G706 -- PathEscpae is sanitizing the url path
		log.Printf("%s %s %s", r.Method, url.PathEscape(r.URL.Path), time.Since(start))
	})
}
