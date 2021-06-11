package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type App struct {
	DB           *pgxpool.Pool
	Config       AppConfig
	MessageQueue *nats.EncodedConn
}

// Fields of a Todo which are not generated.
type NewTodoRequest struct {
	Text string `json:"text"`
}

type UpdateTodoRequest struct {
	Done bool `json:"done"`
}

type Todo struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	Done      bool      `json:"done" db:"done"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func respondJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}

func (app *App) GetImage(w http.ResponseWriter, r *http.Request) {
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

func (app *App) ListTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(r.Context(), `SELECT * FROM todo`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var todos = make([]*Todo, 0)
	pgxscan.ScanAll(&todos, rows)
	respondJSON(w, todos, http.StatusOK)
}

func (app *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var data NewTodoRequest
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if data.Text == "" || len(data.Text) > 140 {
		log.Warn().Str("text", data.Text).Msg("Rejecting due to invalid input")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	db, err := app.DB.Acquire(r.Context())
	if err != nil {
		panic(err)
	}
	defer db.Release()

	var id uuid.UUID
	if err := db.QueryRow(r.Context(), `INSERT INTO todo (text) VALUES ($1) RETURNING id`, data.Text).Scan(&id); err != nil {
		panic(err)
	}
	var todo Todo
	if err := pgxscan.Get(r.Context(), db, &todo, `SELECT * FROM todo WHERE id = $1`, id); err != nil {
		panic(err)
	}

	log.Info().Stringer("id", todo.Id).Str("text", todo.Text).Msg("Created new todo")
	app.MessageQueue.Publish("todos", todo)
	respondJSON(w, todo, http.StatusCreated)
}

func (app *App) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var data UpdateTodoRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	db, err := app.DB.Acquire(r.Context())
	if err != nil {
		panic(err)
	}
	defer db.Release()

	db.Exec(r.Context(), `UPDATE todo SET done = $1 WHERE id = $2`, data.Done, id)
	var todo Todo
	if err := pgxscan.Get(r.Context(), db, &todo, `SELECT * FROM todo WHERE id = $1`, id); err != nil {
		panic(err)
	}
	respondJSON(w, todo, http.StatusOK)
}

// HealthCheck returns 200 OK if the connection to the database is healthy and 503 Service Unavailable otherwise.
func (app *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := app.DB.Ping(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
}
