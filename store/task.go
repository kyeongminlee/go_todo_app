package store

import (
	"context"
	"go_todo_app/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer, id entity.UserID) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
				id
     			, user_id
				, title
				, status
				, created
				, modified
			FROM task
			WHERE user_id = ?;
	`
	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, task *entity.Task) error {
	task.Created = r.Clocker.Now()
	task.Modified = r.Clocker.Now()
	sql := `INSERT INTO task (user_id, title, status, created, modified)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.ExecContext(
		ctx, sql, task.UserID, task.Title, task.Status, task.Created, task.Modified,
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
