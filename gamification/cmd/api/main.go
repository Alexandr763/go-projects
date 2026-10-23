package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gamification/internal/app"
	"gamification/internal/config"
	"gamification/internal/repository/postgres"
	httpapi "gamification/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.LogLevel)
	dbCtx, cancelDB := context.WithTimeout(context.Background(), 3*time.Second)
	db, err := postgres.Open(dbCtx, cfg.DatabaseDSN, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	cancelDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, "database is unavailable")
		os.Exit(1)
	}
	defer db.Close()

	auth := app.New(postgres.New(db), []byte(cfg.JWTSecret), cfg.JWTIssuer, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.Handler(logger, auth),
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown", "err", err)
		}
	}()

	logger.Info("listening", "addr", cfg.HTTPAddr, "env", cfg.AppEnv)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server", "err", err)
		os.Exit(1)
	}
}

func newLogger(level string) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
}
