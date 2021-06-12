package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/strax/dwk-2021/main-app/pkg/protos/pingpong"
)

func mustLookupEnv(key string) string {
	value, ok := os.LookupEnv("MESSAGE")
	if !ok {
		panic(fmt.Errorf("missing environment variable: %v", key))
	}
	return value
}

func main() {
	message := mustLookupEnv("MESSAGE")

	var srv http.Server

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

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

	conn, err := grpc.Dial("pingpong:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	pingpongClient := pingpong.NewPingpongServiceClient(conn)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", middleware.GetReqID(r.Context()))
			next.ServeHTTP(w, r)
		})
	})
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)

	r.Use(hlog.NewHandler(log.Logger))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			newLog := hlog.FromRequest(r).With().Str("requestId", requestId).Logger()
			wrap := hlog.NewHandler(newLog)
			wrap(next).ServeHTTP(w, r)
		})
	})
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		log.Ctx(r.Context()).Trace().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("request")
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, err := pingpongClient.GetStats(r.Context(), &emptypb.Empty{})
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		stats, err := pingpongClient.GetStats(r.Context(), &emptypb.Empty{})
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%v\n%v: %v\nPing / Pongs: %v", message, time.Now().UTC().Format(time.RFC3339), uuid.Must(uuid.NewRandom()), stats.Pings)
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
