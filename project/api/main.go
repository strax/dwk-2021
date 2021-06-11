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

	"github.com/jackc/pgx/v4"
	pgxzerolog "github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	var srv http.Server

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	config := AppConfigFromEnv()
	dbConfig, err := pgxpool.ParseConfig(config.ToPostgresConnectionString())
	if err != nil {
		panic(err)
	}
	dbConfig.LazyConnect = true
	dbConfig.ConnConfig.LogLevel = pgx.LogLevelTrace
	dbConfig.ConnConfig.Logger = pgxzerolog.NewLogger(log.Logger)
	db, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		panic(err)
	}

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
		r.Put("/todos/{id}", app.UpdateTodo)
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
