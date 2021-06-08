package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	var srv http.Server

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	config := AppConfigFromEnv()
	db := sqlx.MustOpen("pgx", config.DBConfig.ToPostgresConnectionString())

	log.Info().Msgf("Service path prefix: %v", config.PathPrefix)

	exit := make(chan struct{})
	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Wait for termination signal
		<-term
		log.Warn().Msg("Shutting down")
		// Shut down the server
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Err(err)
		}
		close(exit)
	}()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(requestIdHeader)
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)

	r.Use(hlog.NewHandler(log.Logger))
	r.Use(requestIdLogger)

	app := App{DB: db, Config: config}

	r.Get("/healthz", app.HealthCheck)

	r.Route(config.PathPrefix, func(r chi.Router) {
		// The rest of the middlewares and routes have the prefix stripped out of the URL path
		r.Use(requestLogger)
		r.Get("/image", app.GetImage)
		r.Post("/todos", app.CreateTodo)
		r.Get("/todos", app.ListTodos)
	})

	srv = http.Server{
		Addr:    "[::]:80",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}

	<-exit
}
