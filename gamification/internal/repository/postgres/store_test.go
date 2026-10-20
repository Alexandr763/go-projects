package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"gamification/internal/app"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordIsStoredAsHash(t *testing.T) {
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
	t.Cleanup(func() {
		db.Close()
	})

	svc := app.NewForTest(New(db), []byte("0123456789abcdef0123456789abcdef"), "gamification", time.Minute, time.Hour, bcrypt.MinCost)
	email := "stage4-" + time.Now().UTC().Format("20060102150405.000000000") + "@example.com"
	session, err := svc.Register(ctx, email, "longenough")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		clean, cancelClean := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelClean()
		if _, err := db.ExecContext(clean, `DELETE FROM refresh_tokens WHERE user_id = $1`, session.User.ID); err != nil {
			t.Errorf("cleanup refresh: %v", err)
		}
		if _, err := db.ExecContext(clean, `DELETE FROM users WHERE id = $1`, session.User.ID); err != nil {
			t.Errorf("cleanup user: %v", err)
		}
	})

	var hash string
	if err := db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id = $1`, session.User.ID).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	if hash == "longenough" || hash == "" {
		t.Fatalf("password stored as %q", hash)
	}
	if _, err := svc.Refresh(ctx, session.RefreshToken); err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(ctx, session.AccessToken, session.RefreshToken); err == nil {
		t.Fatal("old refresh was still valid after rotation")
	}

	var leftover int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM refresh_tokens WHERE user_id = $1 AND revoked_at IS NULL`, session.User.ID).Scan(&leftover); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if leftover != 1 {
		t.Fatalf("live refresh rows = %d", leftover)
	}
}
