package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

var db = make(map[uuid.UUID]Todo)

// Fields of a Todo which are not generated.
type EphemeralTodo struct {
	Text string 	`json:"text"`
}

type Todo struct {
	Id uuid.UUID 	`json:"id"`
	Text string		`json:"text"`
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	for k := range db {
		todos = append(todos, db[k])
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := &EphemeralTodo{}
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
	}
	db[id] = todo
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(todo); err != nil {
		panic(err)
	}
}