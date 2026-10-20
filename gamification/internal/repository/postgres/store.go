package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gamification/internal/app"
	"github.com/jackc/pgx/v5/pgconn"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func Open(ctx context.Context, dsn string, maxOpen, maxIdle int, lifetime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(lifetime)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (app.User, error) {
	var user app.User
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id::text, email, coins, created_at
	`, email, passwordHash).Scan(&user.ID, &user.Email, &user.Coins, &user.CreatedAt)
	if err != nil {
		if isUnique(err) {
			return app.User{}, app.ErrEmailTaken
		}
		return app.User{}, err
	}
	return user, nil
}

func (s *Store) UserByEmail(ctx context.Context, email string) (app.User, error) {
	return s.user(ctx, `SELECT id::text, email, password_hash, coins, created_at FROM users WHERE email = $1`, email)
}

func (s *Store) UserByID(ctx context.Context, id string) (app.User, error) {
	return s.user(ctx, `SELECT id::text, email, password_hash, coins, created_at FROM users WHERE id = $1`, id)
}

func (s *Store) user(ctx context.Context, query, arg string) (app.User, error) {
	var user app.User
	err := s.db.QueryRowContext(ctx, query, arg).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Coins, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return app.User{}, app.ErrNotFound
	}
	if err != nil {
		return app.User{}, err
	}
	return user, nil
}

func (s *Store) SaveRefresh(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (s *Store) UseRefresh(ctx context.Context, tokenHash, nextHash string, expiresAt time.Time) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	var userID string
	err = tx.QueryRowContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > now()
		RETURNING user_id::text
	`, tokenHash).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", app.ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, nextHash, expiresAt); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Store) RevokeRefresh(ctx context.Context, userID, tokenHash string) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE token_hash = $1
		  AND user_id = $2
		  AND revoked_at IS NULL
		  AND expires_at > now()
	`, tokenHash, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return app.ErrUnauthorized
	}
	return nil
}

func (s *Store) CreateTask(ctx context.Context, userID, title, description string, reward int) (app.Task, error) {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO tasks (user_id, title, description, reward_coins)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
	`, userID, title, description, reward)
	return scanTask(row)
}

func (s *Store) ListTasks(ctx context.Context, userID, status string, limit, offset int) ([]app.Task, error) {
	var done any
	switch status {
	case "open":
		done = false
	case "done":
		done = true
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		  AND ($2::boolean IS NULL OR is_done = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, done, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]app.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Store) TaskByID(ctx context.Context, userID, taskID string) (app.Task, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`, taskID, userID)
	task, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return app.Task{}, app.ErrTaskNotFound
	}
	return task, err
}

func (s *Store) UpdateTask(ctx context.Context, userID, taskID string, patch app.TaskPatch) (app.Task, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.Task{}, err
	}
	defer func() { _ = tx.Rollback() }()

	current, err := scanTask(tx.QueryRowContext(ctx, `
		SELECT id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, taskID, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return app.Task{}, app.ErrTaskNotFound
	}
	if err != nil {
		return app.Task{}, err
	}
	if patch.RewardCoins != nil && current.IsDone && *patch.RewardCoins != current.RewardCoins {
		return app.Task{}, app.ErrRewardLocked
	}
	if patch.Title != nil {
		current.Title = *patch.Title
	}
	if patch.Description != nil {
		current.Description = *patch.Description
	}
	if patch.RewardCoins != nil {
		current.RewardCoins = *patch.RewardCoins
	}
	updated, err := scanTask(tx.QueryRowContext(ctx, `
		UPDATE tasks
		SET title = $3, description = $4, reward_coins = $5, updated_at = now()
		WHERE id = $1 AND user_id = $2
		RETURNING id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
	`, taskID, userID, current.Title, current.Description, current.RewardCoins))
	if err != nil {
		return app.Task{}, err
	}
	if err := tx.Commit(); err != nil {
		return app.Task{}, err
	}
	return updated, nil
}

func (s *Store) DeleteTask(ctx context.Context, userID, taskID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var done bool
	err = tx.QueryRowContext(ctx, `
		SELECT is_done FROM tasks WHERE id = $1 AND user_id = $2 FOR UPDATE
	`, taskID, userID).Scan(&done)
	if errors.Is(err, sql.ErrNoRows) {
		return app.ErrTaskNotFound
	}
	if err != nil {
		return err
	}
	if done {
		return app.ErrTaskDone
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1 AND user_id = $2 AND is_done = false`, taskID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return app.ErrTaskNotFound
	}
	return tx.Commit()
}

func (s *Store) CompleteTask(ctx context.Context, userID, taskID string) (app.Task, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return app.Task{}, err
	}
	defer func() { _ = tx.Rollback() }()

	task, err := scanTask(tx.QueryRowContext(ctx, `
		SELECT id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, taskID, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return app.Task{}, app.ErrTaskNotFound
	}
	if err != nil {
		return app.Task{}, err
	}
	if task.IsDone {
		return task, nil
	}

	task, err = scanTask(tx.QueryRowContext(ctx, `
		UPDATE tasks
		SET is_done = true, completed_at = now(), updated_at = now()
		WHERE id = $1 AND user_id = $2
		RETURNING id::text, user_id::text, title, description, reward_coins, is_done, completed_at, created_at, updated_at
	`, taskID, userID))
	if err != nil {
		return app.Task{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE users SET coins = coins + $2, updated_at = now() WHERE id = $1
	`, userID, task.RewardCoins); err != nil {
		return app.Task{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO coin_transactions (user_id, task_id, amount, type)
		VALUES ($1, $2, $3, 'task_reward')
	`, userID, taskID, task.RewardCoins); err != nil {
		if isUnique(err) {
			_ = tx.Rollback()
			return s.awardAlreadyRecorded(ctx, userID, taskID)
		}
		return app.Task{}, err
	}
	if err := tx.Commit(); err != nil {
		return app.Task{}, err
	}
	return task, nil
}

func (s *Store) awardAlreadyRecorded(ctx context.Context, userID, taskID string) (app.Task, error) {
	task, err := s.TaskByID(ctx, userID, taskID)
	if err != nil {
		return app.Task{}, err
	}
	if task.IsDone {
		return task, nil
	}
	return app.Task{}, errAwardRejected
}

func (s *Store) ListTransactions(ctx context.Context, userID string, limit, offset int) ([]app.Transaction, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, user_id::text, task_id::text, amount, type, created_at
		FROM coin_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]app.Transaction, 0)
	for rows.Next() {
		var item app.Transaction
		if err := rows.Scan(&item.ID, &item.UserID, &item.TaskID, &item.Amount, &item.Type, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

var errAwardRejected = errors.New("award rejected")

func scanTask(row interface{ Scan(dest ...any) error }) (app.Task, error) {
	var task app.Task
	var completed sql.NullTime
	err := row.Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.RewardCoins,
		&task.IsDone,
		&completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return app.Task{}, err
	}
	if completed.Valid {
		task.CompletedAt = &completed.Time
	}
	return task, nil
}

func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
