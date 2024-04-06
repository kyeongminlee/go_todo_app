package handler

import (
	"go_todo_app/entity"
	"go_todo_app/store"
	"net/http"
)

type ListTask struct {
	Store *store.TaskStore
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	tasks := lt.Store.All()
	response := []task{}
	for _, t := range tasks {
		response = append(response, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, writer, response, http.StatusOK)
}
