package main

import (
	"bufio"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func image(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("/mnt/volume1/picsum-400-400.webp")
	if err != nil {
		panic(err)
	}
	rdr := bufio.NewReader(f)
	w.Header().Set("Content-Type", "image/webp")
	_, err = rdr.WriteTo(w)
	if err != nil {
		panic(err)
	}
}

func main() {
	var srv http.Server

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	exit := make(chan struct{})
	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Wait for termination signal
		<- term
		log.Warn().Msg("shutting down")
		// Shut down the server
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Err(err)
		}
		close(exit)
	}()
	
	r := chi.NewRouter()
	r.Use(hlog.NewHandler(log.Logger))
	r.Use(hlog.AccessHandler(func (r *http.Request, status, size int, duration time.Duration) {
		log.Ctx(r.Context()).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("request")
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	r.Get("/image", image)
	r.Post("/todos", CreateTodo)
	r.Get("/todos", ListTodos)

	srv = http.Server{
		Addr:    "[::]:80",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}

	<- exit
}