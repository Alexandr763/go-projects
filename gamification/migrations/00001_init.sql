-- +goose Up

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text NOT NULL,
    password_hash text NOT NULL,
    coins integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_email_normalized CHECK (email = lower(btrim(email)) AND email <> ''),
    CONSTRAINT users_password_hash_present CHECK (password_hash <> ''),
    CONSTRAINT users_coins_nonnegative CHECK (coins >= 0)
);

CREATE TABLE tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    reward_coins integer NOT NULL,
    is_done boolean NOT NULL DEFAULT false,
    completed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT tasks_title_present CHECK (btrim(title) <> ''),
    CONSTRAINT tasks_reward_range CHECK (reward_coins >= 0 AND reward_coins <= 10000),
    CONSTRAINT tasks_done_matches_time CHECK (
        (is_done = false AND completed_at IS NULL)
        OR (is_done = true AND completed_at IS NOT NULL)
    ),
    CONSTRAINT tasks_user_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT tasks_id_user_key UNIQUE (id, user_id)
);

CREATE INDEX tasks_user_created_at_idx ON tasks (user_id, created_at DESC);

CREATE TABLE coin_transactions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    task_id uuid NOT NULL,
    amount integer NOT NULL,
    type text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT coin_transactions_amount_nonnegative CHECK (amount >= 0),
    CONSTRAINT coin_transactions_type_reward CHECK (type = 'task_reward'),
    CONSTRAINT coin_transactions_user_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT coin_transactions_task_owner_fkey FOREIGN KEY (task_id, user_id) REFERENCES tasks (id, user_id) ON DELETE RESTRICT,
    CONSTRAINT coin_transactions_task_key UNIQUE (task_id)
);

CREATE INDEX coin_transactions_user_created_at_idx ON coin_transactions (user_id, created_at DESC);

-- +goose Down

DROP TABLE coin_transactions;
DROP TABLE tasks;
DROP TABLE users;
