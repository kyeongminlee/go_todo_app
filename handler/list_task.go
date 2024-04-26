package handler

import (
	"github.com/jmoiron/sqlx"
	"go_todo_app/entity"
	"go_todo_app/store"
	"net/http"
)

type ListTask struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	tasks, err := lt.Repo.ListTasks(ctx, lt.DB)
	if err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	var response []task
	for _, t := range tasks {
		response = append(response, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, writer, response, http.StatusOK)
}
