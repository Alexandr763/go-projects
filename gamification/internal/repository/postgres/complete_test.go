package postgres

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"gamification/internal/app"
	"golang.org/x/crypto/bcrypt"
)

func TestOneCompletionOneAward(t *testing.T) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		t.Skip("DATABASE_DSN is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := Open(ctx, dsn, 4, 2, time.Minute)
	if err != nil {
		t.Skip("database is unavailable")
	}
	t.Cleanup(func() { db.Close() })

	svc := app.NewForTest(New(db), []byte("0123456789abcdef0123456789abcdef"), "gamification", time.Minute, time.Hour, bcrypt.MinCost)
	stamp := time.Now().UTC().Format("20060102150405.000000000")
	owner, err := svc.Register(ctx, "stage6-"+stamp+"@example.com", "longenough")
	if err != nil {
		t.Fatal(err)
	}
	other, err := svc.Register(ctx, "stage6-other-"+stamp+"@example.com", "longenough")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		clean, cancelClean := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelClean()
		for _, id := range []string{owner.User.ID, other.User.ID} {
			if _, err := db.ExecContext(clean, `DELETE FROM coin_transactions WHERE user_id = $1`, id); err != nil {
				t.Errorf("cleanup history: %v", err)
			}
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

	task, err := svc.CreateTask(ctx, owner.User.ID, "ship", "", 15)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.CompleteTask(ctx, owner.User.ID, task.ID)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.CompleteTask(ctx, owner.User.ID, task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteTask(ctx, other.User.ID, task.ID); !errors.Is(err, app.ErrTaskNotFound) {
		t.Fatalf("foreign complete: %v", err)
	}

	var coins, rows, total int
	if err := db.QueryRowContext(ctx, `SELECT coins FROM users WHERE id = $1`, owner.User.ID).Scan(&coins); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, `SELECT count(*), coalesce(sum(amount), 0) FROM coin_transactions WHERE user_id = $1`, owner.User.ID).Scan(&rows, &total); err != nil {
		t.Fatal(err)
	}
	if coins != 15 || rows != 1 || total != 15 {
		t.Fatalf("coins=%d rows=%d total=%d", coins, rows, total)
	}

	if _, err := db.ExecContext(ctx, `UPDATE tasks SET is_done = false, completed_at = NULL WHERE id = $1`, task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteTask(ctx, owner.User.ID, task.ID); err == nil {
		t.Fatal("second insert was accepted")
	}
	var done bool
	if err := db.QueryRowContext(ctx, `SELECT is_done, coins FROM tasks JOIN users ON users.id = tasks.user_id WHERE tasks.id = $1`, task.ID).Scan(&done, &coins); err != nil {
		t.Fatal(err)
	}
	if done || coins != 15 {
		t.Fatalf("after rejected insert done=%v coins=%d", done, coins)
	}
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM coin_transactions WHERE task_id = $1`, task.ID).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows != 1 {
		t.Fatalf("history rows = %d", rows)
	}
}
