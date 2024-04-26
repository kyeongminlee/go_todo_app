package main

import (
	"context"
	"go_todo_app/clock"
	"go_todo_app/config"
	"go_todo_app/handler"
	"go_todo_app/service"
	"go_todo_app/store"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Write([]byte(`{"status": "ok"}`))
	})
	validator := validator.New()

	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		log.Fatalf("db setting error: %v", err)
		return nil, cleanup, err
	}
	repository := store.Repository{
		Clocker: clock.RealClocker{},
	}

	addTask := &handler.AddTask{
		Service: &service.AddTask{
			DB:   db,
			Repo: &repository,
		},
		Validator: validator,
	}
	mux.Post("/tasks", addTask.ServeHTTP)

	listTask := &handler.ListTask{
		Service: &service.ListTask{
			DB:   db,
			Repo: &repository,
		},
	}
	mux.Get("/tasks", listTask.ServeHTTP)

	return mux, cleanup, nil
}
