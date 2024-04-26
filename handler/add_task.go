package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"net/http"
)

type AddTask struct {
	Service   AddTaskService
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

	task, err := at.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	response := struct {
		ID int `json:"id"`
	}{ID: int(task.ID)}
	RespondJSON(ctx, writer, response, http.StatusOK)
}
