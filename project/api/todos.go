package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var db = make(map[uuid.UUID]Todo)

// Fields of a Todo which are not generated.
type NewTodo struct {
	Text string 	`json:"text"`
}

type Todo struct {
	Id uuid.UUID 		`json:"id"`
	Text string			`json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func respondJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	todos := make([]Todo, 0)
	for k := range db {
		todos = append(todos, db[k])
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
	if data.Text == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	todo := Todo {
		Id: id,
		Text: data.Text,
		CreatedAt: time.Now(),
	}
	db[id] = todo
	respondJSON(w, todo, http.StatusCreated)
}