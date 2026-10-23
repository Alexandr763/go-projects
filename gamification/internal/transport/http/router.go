package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"gamification/internal/app"
	"gamification/internal/ratelimit"
	"gamification/internal/transport/http/apigen"
)

type requestIDKey struct{}

func Handler(logger *slog.Logger, auth *app.Service) http.Handler {
	mux := http.NewServeMux()
	apigen.HandlerWithOptions(&API{auth: auth, limit: ratelimit.New(10, time.Minute)}, apigen.StdHTTPServerOptions{
		BaseRouter: mux,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Info("invalid parameter", "request_id", requestID(r), "err", err.Error())
			writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "invalid parameter")
		},
	})
	registerOtherMethods(mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, apigen.ErrorBodyErrorCodeNotFound, "route not found")
	})
	return withRequestID(logger, mux)
}

func withRequestID(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if id == "" || len(id) > 200 || strings.ContainsAny(id, "\r\n") {
			id = newID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.Info("request", "request_id", id, "method", r.Method, "path", r.URL.Path)
	})
}

func requestID(r *http.Request) string {
	id, _ := r.Context().Value(requestIDKey{}).(string)
	return id
}

func newID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(buf)
}
