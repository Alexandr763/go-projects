package app

import (
	"context"
	"strings"
	"time"
)

type Task struct {
	ID          string
	UserID      string
	Title       string
	Description string
	RewardCoins int
	IsDone      bool
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskPatch struct {
	Title       *string
	Description *string
	RewardCoins *int
}

type Transaction struct {
	ID        string
	UserID    string
	TaskID    string
	Amount    int
	Type      string
	CreatedAt time.Time
}

func (s *Service) CreateTask(ctx context.Context, userID, title, description string, reward int) (Task, error) {
	return s.store.CreateTask(ctx, userID, strings.TrimSpace(title), description, reward)
}

func (s *Service) ListTasks(ctx context.Context, userID, status string, limit, offset int) ([]Task, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.store.ListTasks(ctx, userID, status, limit, offset)
}

func (s *Service) Task(ctx context.Context, userID, taskID string) (Task, error) {
	return s.store.TaskByID(ctx, userID, taskID)
}

func (s *Service) UpdateTask(ctx context.Context, userID, taskID string, patch TaskPatch) (Task, error) {
	if patch.Title != nil {
		title := strings.TrimSpace(*patch.Title)
		patch.Title = &title
	}
	return s.store.UpdateTask(ctx, userID, taskID, patch)
}

func (s *Service) DeleteTask(ctx context.Context, userID, taskID string) error {
	return s.store.DeleteTask(ctx, userID, taskID)
}

func (s *Service) CompleteTask(ctx context.Context, userID, taskID string) (Task, error) {
	return s.store.CompleteTask(ctx, userID, taskID)
}

func (s *Service) ListTransactions(ctx context.Context, userID string, limit, offset int) ([]Transaction, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.store.ListTransactions(ctx, userID, limit, offset)
}
