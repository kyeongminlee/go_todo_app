package handler

import (
	"encoding/json"
	"go_todo_app/entity"
	"go_todo_app/store"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(request.Body).Decode(&b); err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	task := &entity.Task{
		Title:   b.Title,
		Status:  "todo",
		Created: time.Now(),
	}
	id, err := store.Tasks.Add(task)
	if err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	response := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, writer, response, http.StatusOK)
}
