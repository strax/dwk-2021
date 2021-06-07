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

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(db *sqlx.DB) {
	start := time.Now()
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != migrate.ErrNoChange {
		panic(err)
	}
	log.Info().Msgf("Ran migrations in %v ms", time.Since(start).Milliseconds())
}

func main() {
	var srv http.Server

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	config := AppConfigFromEnv()
	db := sqlx.MustConnect("pgx", config.DBConfig.ToPostgresConnectionString())

	log.Info().Msgf("Service path prefix: %v", config.PathPrefix)

	runMigrations(db)

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

	r.Use(middleware.Heartbeat("/healthz"))

	// The rest of the middlewares and routes have the prefix stripped out of the URL path
	r.Use(stripPrefix(config.PathPrefix))

	r.Use(hlog.NewHandler(log.Logger))
	r.Use(requestIdLogger)
	r.Use(requestLogger)

	app := App{DB: db, Config: config}

	r.Get("/image", app.GetImage)
	r.Post("/todos", app.CreateTodo)
	r.Get("/todos", app.ListTodos)

	srv = http.Server{
		Addr:    "[::]:80",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}

	<-exit
}
