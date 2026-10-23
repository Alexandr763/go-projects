package httpapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gamification/internal/app"
	"gamification/internal/repository/memory"
	"gamification/internal/transport/http/apigen"
	"golang.org/x/crypto/bcrypt"
)

func TestStage3Responses(t *testing.T) {
	handler := testHandler(t)
	taskID := "11111111-1111-1111-1111-111111111111"
	session := register(t, handler, "owner@example.com", "longenough")
	valid := "Bearer " + session.Access

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		header     string
		requestID  string
		wantStatus int
		wantCode   string
	}{
		{name: "bad json", method: http.MethodPost, path: "/api/v1/auth/register", body: "{", wantStatus: http.StatusBadRequest, wantCode: "validation_error"},
		{name: "short password", method: http.MethodPost, path: "/api/v1/auth/register", body: `{"email":"a@b.c","password":"short"}`, wantStatus: http.StatusBadRequest, wantCode: "validation_error"},
		{name: "missing token", method: http.MethodGet, path: "/api/v1/tasks", wantStatus: http.StatusUnauthorized, wantCode: "unauthorized"},
		{name: "malformed token", method: http.MethodGet, path: "/api/v1/tasks", header: "Bearer abc", wantStatus: http.StatusUnauthorized, wantCode: "unauthorized"},
		{name: "shaped token", method: http.MethodGet, path: "/api/v1/tasks", header: "Bearer aaa.bbb.ccc", wantStatus: http.StatusUnauthorized, wantCode: "unauthorized"},
		{name: "unknown route", method: http.MethodGet, path: "/missing", wantStatus: http.StatusNotFound, wantCode: "not_found"},
		{name: "wrong method", method: http.MethodPut, path: "/api/v1/tasks", wantStatus: http.StatusMethodNotAllowed, wantCode: "method_not_allowed"},
		{name: "bad reward", method: http.MethodPost, path: "/api/v1/tasks", header: valid, body: `{"title":"t","reward_coins":10001}`, wantStatus: http.StatusBadRequest, wantCode: "validation_error"},
		{name: "bad limit", method: http.MethodGet, path: "/api/v1/tasks?limit=0", header: valid, wantStatus: http.StatusBadRequest, wantCode: "validation_error"},
		{name: "bad task id", method: http.MethodGet, path: "/api/v1/tasks/not-a-uuid", header: "Bearer aaa.bbb.ccc", wantStatus: http.StatusBadRequest, wantCode: "validation_error"},
		{name: "task missing", method: http.MethodGet, path: "/api/v1/tasks/" + taskID, header: valid, wantStatus: http.StatusNotFound, wantCode: "not_found"},
		{name: "complete missing task", method: http.MethodPost, path: "/api/v1/tasks/" + taskID + "/complete", header: valid, wantStatus: http.StatusNotFound, wantCode: "not_found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := call(handler, tt.method, tt.path, tt.body, tt.header, "req-1")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}
			if rec.Header().Get("X-Request-ID") != "req-1" {
				t.Fatalf("request id = %q", rec.Header().Get("X-Request-ID"))
			}
			var payload apigen.ErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("json: %v body=%s", err, rec.Body.String())
			}
			if payload.Error.Code != apigen.ErrorBodyErrorCode(tt.wantCode) {
				t.Fatalf("code = %s", payload.Error.Code)
			}
			if strings.Contains(rec.Body.String(), "short") || strings.Contains(rec.Body.String(), "postgres") {
				t.Fatalf("response leaks input: %s", rec.Body.String())
			}
		})
	}
}

func TestHealth(t *testing.T) {
	handler := testHandler(t)

	live := call(handler, http.MethodGet, "/health/live", "", "", "")
	if live.Code != http.StatusOK {
		t.Fatalf("live status = %d", live.Code)
	}
	if live.Header().Get("X-Request-ID") == "" {
		t.Fatal("live request id is empty")
	}

	ready := call(handler, http.MethodGet, "/health/ready", "", "", "")
	if ready.Code != http.StatusServiceUnavailable {
		t.Fatalf("ready status = %d", ready.Code)
	}
}

func call(handler http.Handler, method, path, body, auth, requestID string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

const testSecret = "0123456789abcdef0123456789abcdef"

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	svc := app.NewForTest(memory.New(), []byte(testSecret), "gamification", 15*time.Minute, time.Hour, bcrypt.MinCost)
	return Handler(discardLogger(), svc)
}

func register(t *testing.T, handler http.Handler, email, password string) apigen.Session {
	t.Helper()
	rec := call(handler, http.MethodPost, "/api/v1/auth/register", `{"email":"`+email+`","password":"`+password+`"}`, "", "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("register status = %d body = %s", rec.Code, rec.Body.String())
	}
	var session apigen.Session
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	return session
}
