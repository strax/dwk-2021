package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func stripPrefix(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}

var requestLogger = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
	log.Ctx(r.Context()).Trace().
		Str("method", r.Method).
		Stringer("url", r.URL).
		Int("status", status).
		Int("size", size).
		Dur("duration", duration).
		Msg("request")
})

var requestIdLogger = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := middleware.GetReqID(r.Context())
		newLog := hlog.FromRequest(r).With().Str("requestId", requestId).Logger()
		wrap := hlog.NewHandler(newLog)
		wrap(next).ServeHTTP(w, r)
	})
}

var requestIdHeader = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", middleware.GetReqID(r.Context()))
		next.ServeHTTP(w, r)
	})
}
