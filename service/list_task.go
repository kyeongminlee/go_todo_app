package service

import (
	"context"
	"fmt"
	"go_todo_app/auth"
	"go_todo_app/entity"
	"go_todo_app/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	tasks, err := l.Repo.ListTasks(ctx, l.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	return tasks, nil
}
