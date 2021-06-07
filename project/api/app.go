package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

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

func GetImage(w http.ResponseWriter, r *http.Request) {
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

func ListTodos(w http.ResponseWriter, r *http.Request) {
	todos := make([]Todo, 0)
	db := r.Context().Value("db").(*sqlx.DB)
	err := db.Select(&todos, "SELECT * FROM todo")
	if err != nil {
		panic(err)
	}
	respondJSON(w, todos, http.StatusOK)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
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
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	todo := Todo{
		Id:        id,
		Text:      data.Text,
		CreatedAt: time.Now(),
	}
	db := r.Context().Value("db").(*sqlx.DB)
	_, err = db.NamedExec("INSERT INTO todo (id, text, created_at) VALUES (:id, :text, :created_at)", todo)
	if err != nil {
		panic(err)
	}
	log.Info().Stringer("id", todo.Id).Str("text", todo.Text).Msg("Created new todo")
	respondJSON(w, todo, http.StatusCreated)
}
