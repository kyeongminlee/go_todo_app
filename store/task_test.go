package store

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"go_todo_app/clock"
	"go_todo_app/entity"
	"go_todo_app/testutil"
	"go_todo_app/testutil/fixture"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()
	user := fixture.User(nil)
	result, err := db.ExecContext(ctx, "INSERT INTO user (name, password, role, created, modified) VALUES (?, ?, ?, ?, ?);",
		user.Name, user.Password, user.Role, user.Created, user.Modified)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	return entity.UserID(id)
}

func preparedTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()
	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)

	//if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
	//	t.Logf("failed to initialize task: %v", err)
	//}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID:   userID,
			Title:    "want task 1",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		{
			UserID:   userID,
			Title:    "want task 2",
			Status:   "done",
			Created:  c.Now(),
			Modified: c.Now(),
		},
	}

	tasks := entity.Tasks{
		wants[0],
		{
			UserID:   otherUserID,
			Title:    "not want task",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		wants[1],
	}

	insertQuery := `
		INSERT INTO task (user_id, title, status, created, modified)
		VALUES 
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?);
	`

	result, err := con.ExecContext(
		ctx, insertQuery,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].Created, tasks[0].Modified,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].Created, tasks[1].Modified,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)

	return userID, wants
}

func TestRepository_ListTasks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wantUserID, wants := preparedTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	clock := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  clock.Now(),
		Modified: clock.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	mock.ExpectExec(
		`INSERT INTO task \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: clock}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
