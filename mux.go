package main

import (
	"go_todo_app/handler"
	"go_todo_app/store"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Write([]byte(`{"status": "ok"}`))
	})
	validator := validator.New()
	mux.Handle("/tasks", &handler.AddTask{Store: store.Tasks, Validator: validator})
	addTask := &handler.AddTask{Store: store.Tasks, Validator: validator}
	mux.Post("/tasks", addTask.ServeHTTP)
	listTask := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", listTask.ServeHTTP)
	return mux
}
