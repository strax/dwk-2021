package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type App struct {
	DB     *sqlx.DB
	Config AppConfig
}

// Fields of a Todo which are not generated.
type NewTodo struct {
	Text string `json:"text" db:"text"`
}

type Todo struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
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
	todos := make([]Todo, 0)
	err := app.DB.Select(&todos, "SELECT * FROM todo")
	if err != nil {
		panic(err)
	}
	respondJSON(w, todos, http.StatusOK)
}

func (app *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := &NewTodo{}
	if err := decoder.Decode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if data.Text == "" || len(data.Text) > 140 {
		log.Warn().Str("text", data.Text).Msg("Rejecting due to invalid input")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id := uuid.Must(uuid.NewRandom())
	todo := Todo{
		Id:        id,
		Text:      data.Text,
		CreatedAt: time.Now(),
	}
	_, err := app.DB.NamedExec("INSERT INTO todo (id, text, created_at) VALUES (:id, :text, :created_at)", todo)
	if err != nil {
		panic(err)
	}
	log.Info().Stringer("id", todo.Id).Str("text", todo.Text).Msg("Created new todo")
	respondJSON(w, todo, http.StatusCreated)
}

// HealthCheck returns 200 OK if the connection to the database is healthy and 503 Service Unavailable otherwise.
func (app *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := app.DB.Ping()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
}
