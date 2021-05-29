package main

import (
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

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(db *sqlx.DB) {
	log.Info().Msg("Running migrations")
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
}

func main() {
	var srv http.Server

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	dbconfig := DBConfigFromEnv()
	db := sqlx.MustConnect("pgx", dbconfig.ToPostgresConnectionString())

	runMigrations(db)

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
	r.Use(func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Get("/image", GetImage)
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