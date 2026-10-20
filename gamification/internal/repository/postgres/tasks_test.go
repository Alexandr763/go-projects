package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"gamification/internal/app"
	"golang.org/x/crypto/bcrypt"
)

func TestTaskBelongsToItsUser(t *testing.T) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		t.Skip("DATABASE_DSN is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := Open(ctx, dsn, 2, 1, time.Minute)
	if err != nil {
		t.Skip("database is unavailable")
	}
	t.Cleanup(func() { db.Close() })

	svc := app.NewForTest(New(db), []byte("0123456789abcdef0123456789abcdef"), "gamification", time.Minute, time.Hour, bcrypt.MinCost)
	stamp := time.Now().UTC().Format("20060102150405.000000000")
	owner, err := svc.Register(ctx, "stage5-owner-"+stamp+"@example.com", "longenough")
	if err != nil {
		t.Fatal(err)
	}
	other, err := svc.Register(ctx, "stage5-other-"+stamp+"@example.com", "longenough")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		clean, cancelClean := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelClean()
		for _, id := range []string{owner.User.ID, other.User.ID} {
			if _, err := db.ExecContext(clean, `DELETE FROM tasks WHERE user_id = $1`, id); err != nil {
				t.Errorf("cleanup tasks: %v", err)
			}
			if _, err := db.ExecContext(clean, `DELETE FROM refresh_tokens WHERE user_id = $1`, id); err != nil {
				t.Errorf("cleanup refresh: %v", err)
			}
			if _, err := db.ExecContext(clean, `DELETE FROM users WHERE id = $1`, id); err != nil {
				t.Errorf("cleanup user: %v", err)
			}
		}
	})

	title := "x'); DROP TABLE tasks;--"
	task, err := svc.CreateTask(ctx, owner.User.ID, title, "", 5)
	if err != nil {
		t.Fatal(err)
	}
	var stored string
	var tableCount int
	if err := db.QueryRowContext(ctx, `SELECT title FROM tasks WHERE id = $1`, task.ID).Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if stored != title {
		t.Fatalf("stored title = %q", stored)
	}
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM information_schema.tables WHERE table_name = 'tasks'`).Scan(&tableCount); err != nil {
		t.Fatal(err)
	}
	if tableCount != 1 {
		t.Fatal("tasks table is missing")
	}

	if _, err := svc.Task(ctx, other.User.ID, task.ID); !errors.Is(err, app.ErrTaskNotFound) {
		t.Fatalf("foreign read: %v", err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE tasks SET is_done = true, completed_at = now() WHERE id = $1`, task.ID); err != nil {
		t.Fatal(err)
	}
	reward := 6
	if _, err := svc.UpdateTask(ctx, owner.User.ID, task.ID, app.TaskPatch{RewardCoins: &reward}); !errors.Is(err, app.ErrRewardLocked) {
		t.Fatalf("reward change: %v", err)
	}
	if err := svc.DeleteTask(ctx, owner.User.ID, task.ID); !errors.Is(err, app.ErrTaskDone) {
		t.Fatalf("delete done: %v", err)
	}
	var coins int
	if err := db.QueryRowContext(ctx, `SELECT coins FROM users WHERE id = $1`, owner.User.ID).Scan(&coins); err != nil {
		t.Fatal(err)
	}
	if coins != 0 {
		t.Fatalf("coins = %d", coins)
	}

	stopped, cancelStopped := context.WithCancel(context.Background())
	cancelStopped()
	if _, err := svc.ListTasks(stopped, owner.User.ID, "", 20, 0); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled list: %v", err)
	}
}
