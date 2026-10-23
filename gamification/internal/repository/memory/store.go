package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"gamification/internal/app"
	"github.com/google/uuid"
)

type Store struct {
	mu      sync.Mutex
	users   map[string]app.User
	byEmail map[string]string
	refresh map[string]refreshRow
	tasks   map[string]app.Task
	history []app.Transaction
	awarded map[string]struct{}
	last    time.Time
}

type refreshRow struct {
	userID    string
	expiresAt time.Time
	revoked   bool
}

func New() *Store {
	return &Store{
		users:   map[string]app.User{},
		byEmail: map[string]string{},
		refresh: map[string]refreshRow{},
		tasks:   map[string]app.Task{},
		history: []app.Transaction{},
		awarded: map[string]struct{}{},
	}
}

func (s *Store) CreateUser(_ context.Context, email, passwordHash string) (app.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byEmail[email]; ok {
		return app.User{}, app.ErrEmailTaken
	}
	user := app.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		Coins:        0,
		CreatedAt:    time.Now().UTC(),
	}
	s.users[user.ID] = user
	s.byEmail[email] = user.ID
	return user, nil
}

func (s *Store) UserByEmail(_ context.Context, email string) (app.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.byEmail[email]
	if !ok {
		return app.User{}, app.ErrNotFound
	}
	return s.users[id], nil
}

func (s *Store) UserByID(_ context.Context, id string) (app.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[id]
	if !ok {
		return app.User{}, app.ErrNotFound
	}
	return user, nil
}

func (s *Store) SaveRefresh(_ context.Context, userID, tokenHash string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refresh[tokenHash] = refreshRow{userID: userID, expiresAt: expiresAt}
	return nil
}

func (s *Store) UseRefresh(_ context.Context, tokenHash, nextHash string, expiresAt time.Time) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row, ok := s.refresh[tokenHash]
	if !ok || row.revoked || !row.expiresAt.After(time.Now()) {
		return "", app.ErrUnauthorized
	}
	row.revoked = true
	s.refresh[tokenHash] = row
	s.refresh[nextHash] = refreshRow{userID: row.userID, expiresAt: expiresAt}
	return row.userID, nil
}

func (s *Store) RevokeRefresh(_ context.Context, userID, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	row, ok := s.refresh[tokenHash]
	if !ok || row.revoked || row.userID != userID {
		return app.ErrUnauthorized
	}
	row.revoked = true
	s.refresh[tokenHash] = row
	return nil
}

func (s *Store) CreateTask(ctx context.Context, userID, title, description string, reward int) (app.Task, error) {
	if err := ctx.Err(); err != nil {
		return app.Task{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.nextTime()
	task := app.Task{
		ID:          uuid.NewString(),
		UserID:      userID,
		Title:       title,
		Description: description,
		RewardCoins: reward,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task
	return task, nil
}

func (s *Store) ListTasks(ctx context.Context, userID, status string, limit, offset int) ([]app.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	all := make([]app.Task, 0)
	for _, task := range s.tasks {
		if task.UserID != userID {
			continue
		}
		if status == "open" && task.IsDone {
			continue
		}
		if status == "done" && !task.IsDone {
			continue
		}
		all = append(all, task)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	if offset > len(all) {
		return []app.Task{}, nil
	}
	all = all[offset:]
	if limit < len(all) {
		all = all[:limit]
	}
	return all, nil
}

func (s *Store) TaskByID(ctx context.Context, userID, taskID string) (app.Task, error) {
	if err := ctx.Err(); err != nil {
		return app.Task{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return app.Task{}, app.ErrTaskNotFound
	}
	return task, nil
}

func (s *Store) UpdateTask(ctx context.Context, userID, taskID string, patch app.TaskPatch) (app.Task, error) {
	if err := ctx.Err(); err != nil {
		return app.Task{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return app.Task{}, app.ErrTaskNotFound
	}
	if patch.RewardCoins != nil && task.IsDone && *patch.RewardCoins != task.RewardCoins {
		return app.Task{}, app.ErrRewardLocked
	}
	if patch.Title != nil {
		task.Title = *patch.Title
	}
	if patch.Description != nil {
		task.Description = *patch.Description
	}
	if patch.RewardCoins != nil {
		task.RewardCoins = *patch.RewardCoins
	}
	task.UpdatedAt = s.nextTime()
	s.tasks[taskID] = task
	return task, nil
}

func (s *Store) DeleteTask(ctx context.Context, userID, taskID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return app.ErrTaskNotFound
	}
	if task.IsDone {
		return app.ErrTaskDone
	}
	delete(s.tasks, taskID)
	return nil
}

func (s *Store) CompleteTask(ctx context.Context, userID, taskID string) (app.Task, error) {
	if err := ctx.Err(); err != nil {
		return app.Task{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return app.Task{}, app.ErrTaskNotFound
	}
	if task.IsDone {
		return task, nil
	}
	if _, taken := s.awarded[taskID]; taken {
		return app.Task{}, errAwardRejected
	}
	now := s.nextTime()
	task.IsDone = true
	task.CompletedAt = &now
	task.UpdatedAt = now
	s.tasks[taskID] = task
	if user, ok := s.users[userID]; ok {
		user.Coins += task.RewardCoins
		s.users[userID] = user
	}
	s.awarded[taskID] = struct{}{}
	s.history = append(s.history, app.Transaction{
		ID:        uuid.NewString(),
		UserID:    userID,
		TaskID:    taskID,
		Amount:    task.RewardCoins,
		Type:      "task_reward",
		CreatedAt: now,
	})
	return task, nil
}

func (s *Store) ListTransactions(ctx context.Context, userID string, limit, offset int) ([]app.Transaction, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	rows := make([]app.Transaction, 0)
	for _, row := range s.history {
		if row.UserID == userID {
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CreatedAt.After(rows[j].CreatedAt)
	})
	if offset > len(rows) {
		return []app.Transaction{}, nil
	}
	rows = rows[offset:]
	if limit < len(rows) {
		rows = rows[:limit]
	}
	return rows, nil
}

var errAwardRejected = errors.New("award rejected")

func (s *Store) MarkDone(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[taskID]
	if !ok {
		return
	}
	now := s.nextTime()
	task.IsDone = true
	task.CompletedAt = &now
	task.UpdatedAt = now
	s.tasks[taskID] = task
}

func (s *Store) nextTime() time.Time {
	now := time.Now().UTC()
	if !now.After(s.last) {
		now = s.last.Add(time.Nanosecond)
	}
	s.last = now
	return now
}
