package httpapi

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"gamification/internal/app"
	"gamification/internal/ratelimit"
	"gamification/internal/transport/http/apigen"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type API struct {
	auth  *app.Service
	limit *ratelimit.Limiter
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var body apigen.Credentials
	if !readJSON(w, r, &body) {
		return
	}
	if !validCredentials(w, body) {
		return
	}
	if !a.allow(w, r, "register") {
		return
	}
	session, err := a.auth.Register(r.Context(), string(body.Email), body.Password)
	a.writeSession(w, http.StatusCreated, session, err)
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var body apigen.Credentials
	if !readJSON(w, r, &body) {
		return
	}
	if !validCredentials(w, body) {
		return
	}
	if !a.allow(w, r, "login") {
		return
	}
	session, err := a.auth.Login(r.Context(), string(body.Email), body.Password)
	a.writeSession(w, http.StatusOK, session, err)
}

func (a *API) Refresh(w http.ResponseWriter, r *http.Request) {
	var body apigen.RefreshToken
	if !readJSON(w, r, &body) {
		return
	}
	if strings.TrimSpace(body.Refresh) == "" {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "refresh is required")
		return
	}
	session, err := a.auth.Refresh(r.Context(), body.Refresh)
	a.writeSession(w, http.StatusOK, session, err)
}

func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := bearer(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, apigen.ErrorBodyErrorCodeUnauthorized, "missing or invalid token")
		return
	}
	var body apigen.RefreshToken
	if !readJSON(w, r, &body) {
		return
	}
	if strings.TrimSpace(body.Refresh) == "" {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "refresh is required")
		return
	}
	if err := a.auth.Logout(r.Context(), token, body.Refresh); err != nil {
		a.writeAuthError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) ListTasks(w http.ResponseWriter, r *http.Request, params apigen.ListTasksParams) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	if params.Limit != nil && (*params.Limit < 1 || *params.Limit > 100) {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "limit is out of range")
		return
	}
	if params.Offset != nil && *params.Offset < 0 {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "offset is out of range")
		return
	}
	status := ""
	if params.Status != nil {
		status = string(*params.Status)
	}
	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}
	tasks, err := a.auth.ListTasks(r.Context(), userID, status, limit, offset)
	if err != nil {
		a.writeTaskError(w, err)
		return
	}
	body := make([]apigen.Task, 0, len(tasks))
	for _, task := range tasks {
		item, err := taskJSON(task)
		if err != nil {
			writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
			return
		}
		body = append(body, item)
	}
	writeJSON(w, http.StatusOK, apigen.TaskList{Tasks: body})
}

func (a *API) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	var body apigen.CreateTask
	if !readJSON(w, r, &body) {
		return
	}
	if strings.TrimSpace(body.Title) == "" {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "title is required")
		return
	}
	if body.RewardCoins < 0 || body.RewardCoins > 10000 {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "reward_coins is out of range")
		return
	}
	description := ""
	if body.Description != nil {
		description = *body.Description
	}
	task, err := a.auth.CreateTask(r.Context(), userID, body.Title, description, body.RewardCoins)
	if err != nil {
		a.writeTaskError(w, err)
		return
	}
	a.writeTask(w, http.StatusCreated, task)
}

func (a *API) GetTask(w http.ResponseWriter, r *http.Request, id apigen.TaskID) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	task, err := a.auth.Task(r.Context(), userID, id.String())
	if err != nil {
		a.writeTaskError(w, err)
		return
	}
	a.writeTask(w, http.StatusOK, task)
}

func (a *API) UpdateTask(w http.ResponseWriter, r *http.Request, id apigen.TaskID) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	var body apigen.UpdateTask
	if !readJSON(w, r, &body) {
		return
	}
	if body.Title == nil && body.Description == nil && body.RewardCoins == nil {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "no fields to update")
		return
	}
	if body.Title != nil && strings.TrimSpace(*body.Title) == "" {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "title is required")
		return
	}
	if body.RewardCoins != nil && (*body.RewardCoins < 0 || *body.RewardCoins > 10000) {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "reward_coins is out of range")
		return
	}
	task, err := a.auth.UpdateTask(r.Context(), userID, id.String(), app.TaskPatch{
		Title:       body.Title,
		Description: body.Description,
		RewardCoins: body.RewardCoins,
	})
	if err != nil {
		a.writeTaskError(w, err)
		return
	}
	a.writeTask(w, http.StatusOK, task)
}

func (a *API) DeleteTask(w http.ResponseWriter, r *http.Request, id apigen.TaskID) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	if err := a.auth.DeleteTask(r.Context(), userID, id.String()); err != nil {
		a.writeTaskError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) CompleteTask(w http.ResponseWriter, r *http.Request, id apigen.TaskID) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	task, err := a.auth.CompleteTask(r.Context(), userID, id.String())
	if err != nil {
		a.writeTaskError(w, err)
		return
	}
	a.writeTask(w, http.StatusOK, task)
}

func (a *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	token, ok := bearer(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, apigen.ErrorBodyErrorCodeUnauthorized, "missing or invalid token")
		return
	}
	user, err := a.auth.Profile(r.Context(), token)
	if err != nil {
		a.writeAuthError(w, err)
		return
	}
	id, err := uuid.Parse(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, apigen.Profile{
		Id:        openapi_types.UUID(id),
		Email:     user.Email,
		Coins:     user.Coins,
		CreatedAt: user.CreatedAt,
	})
}

func (a *API) ListTransactions(w http.ResponseWriter, r *http.Request, params apigen.ListTransactionsParams) {
	userID, ok := a.requireUser(w, r)
	if !ok {
		return
	}
	if params.Limit != nil && (*params.Limit < 1 || *params.Limit > 100) {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "limit is out of range")
		return
	}
	if params.Offset != nil && *params.Offset < 0 {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "offset is out of range")
		return
	}
	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}
	rows, err := a.auth.ListTransactions(r.Context(), userID, limit, offset)
	if err != nil {
		a.writeTaskError(w, err)
		return
	}
	body := make([]apigen.Transaction, 0, len(rows))
	for _, row := range rows {
		item, err := transactionJSON(row)
		if err != nil {
			writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
			return
		}
		body = append(body, item)
	}
	writeJSON(w, http.StatusOK, apigen.TransactionList{Transactions: body})
}

func (a *API) HealthLive(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, apigen.Health{Status: "live"})
}

func (a *API) HealthReady(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusServiceUnavailable, apigen.ErrorBodyErrorCodeServiceUnavailable, "database is not checked")
}

func (a *API) requireUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	token, ok := bearer(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, apigen.ErrorBodyErrorCodeUnauthorized, "missing or invalid token")
		return "", false
	}
	id, err := a.auth.UserID(token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, apigen.ErrorBodyErrorCodeUnauthorized, "missing or invalid token")
		return "", false
	}
	return id, true
}

func (a *API) allow(w http.ResponseWriter, r *http.Request, route string) bool {
	if a.limit.Allow(route + " " + clientIP(r)) {
		return true
	}
	writeError(w, http.StatusTooManyRequests, apigen.ErrorBodyErrorCodeTooManyRequests, "too many requests")
	return false
}

func (a *API) writeSession(w http.ResponseWriter, status int, session app.Session, err error) {
	if err != nil {
		a.writeAuthError(w, err)
		return
	}
	id, err := uuid.Parse(session.User.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
		return
	}
	writeJSON(w, status, apigen.Session{
		Id:      openapi_types.UUID(id),
		Email:   session.User.Email,
		Coins:   session.User.Coins,
		Access:  session.AccessToken,
		Refresh: session.RefreshToken,
	})
}

func (a *API) writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, app.ErrEmailTaken):
		writeError(w, http.StatusConflict, apigen.ErrorBodyErrorCodeConflict, "email is already taken")
	case errors.Is(err, app.ErrCredentials):
		writeError(w, http.StatusUnauthorized, apigen.ErrorBodyErrorCodeUnauthorized, "invalid email or password")
	case errors.Is(err, app.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, apigen.ErrorBodyErrorCodeUnauthorized, "missing or invalid token")
	default:
		writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
	}
}

func (a *API) writeTask(w http.ResponseWriter, status int, task app.Task) {
	body, err := taskJSON(task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
		return
	}
	writeJSON(w, status, body)
}

func (a *API) writeTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, app.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, apigen.ErrorBodyErrorCodeNotFound, "task not found")
	case errors.Is(err, app.ErrRewardLocked):
		writeError(w, http.StatusConflict, apigen.ErrorBodyErrorCodeConflict, "reward cannot change after completion")
	case errors.Is(err, app.ErrTaskDone):
		writeError(w, http.StatusConflict, apigen.ErrorBodyErrorCodeConflict, "completed task cannot be deleted")
	default:
		writeError(w, http.StatusInternalServerError, apigen.ErrorBodyErrorCodeInternalError, "internal error")
	}
}

func taskJSON(task app.Task) (apigen.Task, error) {
	id, err := uuid.Parse(task.ID)
	if err != nil {
		return apigen.Task{}, err
	}
	return apigen.Task{
		Id:          openapi_types.UUID(id),
		Title:       task.Title,
		Description: task.Description,
		RewardCoins: task.RewardCoins,
		IsDone:      task.IsDone,
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func transactionJSON(row app.Transaction) (apigen.Transaction, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return apigen.Transaction{}, err
	}
	taskID, err := uuid.Parse(row.TaskID)
	if err != nil {
		return apigen.Transaction{}, err
	}
	return apigen.Transaction{
		Id:        openapi_types.UUID(id),
		TaskId:    openapi_types.UUID(taskID),
		Amount:    row.Amount,
		Type:      row.Type,
		CreatedAt: row.CreatedAt,
	}, nil
}

func validCredentials(w http.ResponseWriter, body apigen.Credentials) bool {
	email := strings.ToLower(strings.TrimSpace(string(body.Email)))
	if email == "" || !strings.Contains(email, "@") || strings.ContainsAny(email, " \t") {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "invalid email")
		return false
	}
	if len(body.Password) < 8 {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "password must be at least 8 characters")
		return false
	}
	return true
}

func bearer(r *http.Request) (string, bool) {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	token = strings.TrimSpace(token)
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
