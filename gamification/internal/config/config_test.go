package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	useCleanEnv(t)
	t.Setenv("DATABASE_DSN", "postgres://localhost/gamification")
	t.Setenv("JWT_SECRET", strings.Repeat("s", 32))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q", cfg.HTTPAddr)
	}
	if cfg.HTTPReadHeaderTimeout != 5*time.Second {
		t.Fatalf("header timeout = %s", cfg.HTTPReadHeaderTimeout)
	}
	if cfg.DBMaxOpenConns != 10 || cfg.DBMaxIdleConns != 5 {
		t.Fatalf("pool = %d/%d", cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	}
	if cfg.LogLevel != "info" || cfg.AppEnv != "dev" {
		t.Fatalf("level=%s env=%s", cfg.LogLevel, cfg.AppEnv)
	}
	if cfg.JWTAccessTTL != 15*time.Minute || cfg.JWTIssuer != "gamification" {
		t.Fatalf("jwt access=%s issuer=%s", cfg.JWTAccessTTL, cfg.JWTIssuer)
	}
}

func TestLoadRejectsMissingDSN(t *testing.T) {
	useCleanEnv(t)
	t.Setenv("JWT_SECRET", strings.Repeat("s", 32))

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DATABASE_DSN") {
		t.Fatalf("err = %v", err)
	}
}

func TestLoadRejectsShortSecretWithoutValue(t *testing.T) {
	useCleanEnv(t)
	secret := "too-short"
	t.Setenv("DATABASE_DSN", "postgres://localhost/gamification")
	t.Setenv("JWT_SECRET", secret)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error contains secret: %v", err)
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("err = %v", err)
	}
}

func TestLoadRejectsIdlePoolAboveOpen(t *testing.T) {
	useCleanEnv(t)
	t.Setenv("DATABASE_DSN", "postgres://localhost/gamification")
	t.Setenv("JWT_SECRET", strings.Repeat("s", 32))
	t.Setenv("DB_MAX_OPEN_CONNS", "2")
	t.Setenv("DB_MAX_IDLE_CONNS", "5")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DB_MAX_IDLE_CONNS") {
		t.Fatalf("err = %v", err)
	}
}

func useCleanEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"DATABASE_DSN",
		"JWT_SECRET",
		"HTTP_ADDR",
		"HTTP_READ_HEADER_TIMEOUT",
		"HTTP_READ_TIMEOUT",
		"HTTP_WRITE_TIMEOUT",
		"HTTP_IDLE_TIMEOUT",
		"DB_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS",
		"DB_CONN_MAX_LIFETIME",
		"LOG_LEVEL",
		"APP_ENV",
		"JWT_ACCESS_TTL",
		"JWT_REFRESH_TTL",
		"JWT_ISSUER",
	}
	for _, key := range keys {
		t.Setenv(key, "")
		if err := osUnset(key); err != nil {
			t.Fatal(err)
		}
	}
}

func osUnset(key string) error {
	return unsetenv(key)
}
