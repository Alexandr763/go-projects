package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config — настройки запуска. Секрет и адрес базы читаются из окружения и в журнал не пишутся.
type Config struct {
	HTTPAddr              string
	HTTPReadHeaderTimeout time.Duration
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration

	DatabaseDSN       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	LogLevel string
	AppEnv   string

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	JWTIssuer     string
}

// Load читает окружение. При пустой или кривой настройке возвращает ошибку с именем переменной.
func Load() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.DatabaseDSN, err = required("DATABASE_DSN")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.JWTSecret, err = required("JWT_SECRET")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("config: JWT_SECRET must be at least 32 characters")
	}

	cfg.HTTPAddr, err = withDefault("HTTP_ADDR", ":8080")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.HTTPReadHeaderTimeout, err = positiveDuration("HTTP_READ_HEADER_TIMEOUT", "5s")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.HTTPReadTimeout, err = positiveDuration("HTTP_READ_TIMEOUT", "10s")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.HTTPWriteTimeout, err = positiveDuration("HTTP_WRITE_TIMEOUT", "10s")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.HTTPIdleTimeout, err = positiveDuration("HTTP_IDLE_TIMEOUT", "60s")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	cfg.DBMaxOpenConns, err = positiveInt("DB_MAX_OPEN_CONNS", "10")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.DBMaxIdleConns, err = positiveInt("DB_MAX_IDLE_CONNS", "5")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		return nil, fmt.Errorf("config: DB_MAX_IDLE_CONNS must not exceed DB_MAX_OPEN_CONNS")
	}
	cfg.DBConnMaxLifetime, err = positiveDuration("DB_CONN_MAX_LIFETIME", "30m")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	cfg.LogLevel, err = oneOf("LOG_LEVEL", "info", "debug", "info", "warn", "error")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.AppEnv, err = oneOf("APP_ENV", "dev", "dev", "prod")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	cfg.JWTAccessTTL, err = positiveDuration("JWT_ACCESS_TTL", "15m")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.JWTRefreshTTL, err = positiveDuration("JWT_REFRESH_TTL", "720h")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	cfg.JWTIssuer, err = withDefault("JWT_ISSUER", "gamification")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

func required(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)
	if !ok || value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func withDefault(key, fallback string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s is empty", key)
	}
	return value, nil
}

func positiveDuration(key, fallback string) (time.Duration, error) {
	raw, err := withDefault(key, fallback)
	if err != nil {
		return 0, err
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}
	return parsed, nil
}

func positiveInt(key, fallback string) (int, error) {
	raw, err := withDefault(key, fallback)
	if err != nil {
		return 0, err
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return parsed, nil
}

func oneOf(key, fallback string, allowed ...string) (string, error) {
	raw, err := withDefault(key, fallback)
	if err != nil {
		return "", err
	}
	for _, item := range allowed {
		if raw == item {
			return raw, nil
		}
	}
	return "", fmt.Errorf("%s must be one of %s", key, strings.Join(allowed, ", "))
}
