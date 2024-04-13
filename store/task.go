package store

import (
	"context"
	"go_todo_app/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
				id
				, title
				, status
				, created
				, modified
			FROM task;
	`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, task *entity.Task) error {
	task.Created = r.Clocker.Now()
	task.Modified = r.Clocker.Now()
	sql := `INSERT INTO task (title, status, created, modified)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.ExecContext(
		ctx, sql, task.Title, task.Status, task.Created, task.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = entity.TaskID(id)
	return nil
}
